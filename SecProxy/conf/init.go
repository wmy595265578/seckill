package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	redisPoll  *redis.Pool
	etcdClient *etcd_client.Client
)

func initRedis() (err error) {
	redisPoll = &redis.Pool{
		MaxIdle:     SeckillConf.RedisConf.RedisMaxIdle,
		MaxActive:   SeckillConf.RedisConf.RedisMaxActive,
		IdleTimeout: time.Duration(SeckillConf.RedisConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", SeckillConf.RedisConf.RedisAddr, redis.DialPassword(SeckillConf.RedisConf.RedisPassword))
		},
	}
	conn := redisPoll.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("conn redis failed,err:%v", err)
		return
	}
	return
}

func initEtcd() (err error) {
	cli, err := etcd_client.New(etcd_client.Config{
		Endpoints:   []string{SeckillConf.EtcdConf.EtcdAddr},
		DialTimeout: time.Duration(SeckillConf.EtcdConf.EtcdTimeout) * time.Second,
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
	config["filename"] = SeckillConf.LogPath
	config["level"] = convertLogLevel(SeckillConf.LogLevel)

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
	key := fmt.Sprintf("%s/product", SeckillConf.EtcdConf.EtcdSecKeyPrefix)
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	resp, err := etcdClient.Get(ctx, key)
	if err != nil {
		logs.Error("Get from etcd key[%s] config failed,err:%v", key, err)
		return
	}
	for k, v := range resp.Kvs {
		logs.Debug("key[%v] valud[%v]", k, v)
	}
	cancle()
	return
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
	logs.Info("init sec successful ")

	return
}
