package conf

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/gomodule/redigo/redis"
	"time"
	"seckill/SecProxy/models"
)

var (
	redisPoll  *redis.Pool
	etcdClient *etcd_client.Client
)

//func initRedis() (err error) {
//	redisPoll = &redis.Pool{
//		MaxIdle:    models.SeckillConf.RedisProxy2LayerConf.RedisMaxIdle,
//		MaxActive:   models.SeckillConf.RedisProxy2LayerConf.RedisMaxActive,
//		IdleTimeout: time.Duration(models.SeckillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
//		Dial: func() (redis.Conn, error) {
//			return redis.Dial("tcp", models.SeckillConf.RedisProxy2LayerConf.RedisAddr, redis.DialPassword(models.SeckillConf.RedisProxy2LayerConf.RedisPassword))
//		},
//	}
//	conn := redisPoll.Get()
//	defer conn.Close()
//	_, err = conn.Do("ping")
//	if err != nil {
//		logs.Error("conn redis failed,err:%v", err)
//		return
//	}
//	return
//}

func initEtcd() (err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{models.SeckillConf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(models.SeckillConf.EtcdConf.EtcdTimeout) * time.Second,
	})
	if err != nil {
		logs.Error("connect etcd failed, err:", err)
		return
	}

	etcdClient = cli
	return
}

func convertLogLevel(level string) int {

	switch level {
	case "debug":
		return logs.LevelDebug
	case "warn":
		return logs.LevelWarn
	case "info":
		return logs.LevelInfo
	case "trace":
		return logs.LevelTrace
	}

	return logs.LevelDebug
}

func initLogs() (err error) {
	config := make(map[string]interface{})
	config["filename"] = models.SeckillConf.LogPath
	config["level"] = convertLogLevel(models.SeckillConf.LogLevel)

	configStr, err := json.Marshal(config)
	if err != nil {
		logs.Error("marshal failed, err:", err)
		return
	}

	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
	return
}

func loadSecConf() (err error) {
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.Get(ctx, models.SeckillConf.EtcdConf.EtcdProductKey)
	logs.Info("etcd key[%s]", models.SeckillConf.EtcdConf.EtcdProductKey)
	if err != nil {
		logs.Error("Get from etcd key[%s] config failed,err:%v", models.SeckillConf.EtcdConf.EtcdProductKey, err)
		return
	}
	var SecProductInfo []models.SecInfoConfing
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valud[%v]", k, v)
		err = json.Unmarshal(v.Value, &SecProductInfo)
		if err != nil {
			logs.Error("Unmasrshal etcd config failed ,err :%v", err)
			return
		}
		logs.Debug("sec info conf is [%v]", SecProductInfo)
	}
	cancle()
	updateSecProductInfo(SecProductInfo)
	return
}

func initSecProductWatcher() {
	go watcherSecProductKey(models.SeckillConf.EtcdConf.EtcdProductKey)
}

func watcherSecProductKey(key string) {
	logs.Debug("begin watch key:%s", key)
	for {
		rch := etcdClient.Watch(context.Background(), key)
		var SecProductInfo []models.SecInfoConfing
		var getConfSucc = true
		for wresp := range rch {
			for _, ev := range wresp.Events {
				if ev.Type == mvccpb.DELETE {
					logs.Warn("key[%s] 's config deleted", key)
					continue
				}

				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &SecProductInfo)
					if err != nil {
						logs.Error("key[%s] ,Unmarshal failed,err:%v", ev.Kv.Key, err)
						getConfSucc = false
						continue
					}

				}
				logs.Debug("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
			if getConfSucc {
				logs.Debug("get config from etcd succ, %v", SecProductInfo)
				updateSecProductInfo(SecProductInfo)
			}
		}
	}
}
func updateSecProductInfo(SecProductInfo []models.SecInfoConfing) {
	var tmp map[int]*models.SecInfoConfing = make(map[int]*models.SecInfoConfing, 1024)
	for _, v := range SecProductInfo {
		productInfo := v
		tmp[v.ProductId] = &productInfo
	}
	models.SeckillConf.RWSecKillLock.Lock()
	models.SeckillConf.SecInfoConfMap = tmp
	models.SeckillConf.RWSecKillLock.Unlock()
}

func InitSec() (err error) {
	err = initConfig()
	if err != nil {
		logs.Error("init Config failed err:%v", err)
		panic(err)
		return
	}
	err = initLogs()
	if err != nil {
		logs.Error("init logs config failed ,err:%v", err)
		return
	}

	err = initRedis()
	if err != nil {
		logs.Error("init Redis failed err:%v", err)
		panic(err)
		return
	}
	err = initEtcd()
	if err != nil {
		logs.Error("init Etcd failed err:%v", err)
		panic(err)
		return
	}

	err = loadSecConf()
	if err != nil {
		logs.Error("loadSecConf failed err:%v", err)
		panic(err)
		return
	}

	initSecProductWatcher()
	logs.Info("init sec successful ")

	return
}
