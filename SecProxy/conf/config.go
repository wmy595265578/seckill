package conf

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strings"
	"sync"
)

var (
	SeckillConf = &SecSkillConf{
		SecInfoConfMap:make(map[int]*SecInfoConfing),
	}
)

type SecSkillConf struct {
	RedisConf   RedisConfing
	EtcdConf    EtcdConfig
	LogPath     string
	LogLevel    string
	SecInfoConfMap map[int]*SecInfoConfing
	RWSecKillLock     sync.RWMutex

}

type EtcdConfig struct {
	EtcdAddr         string
	EtcdTimeout      int
	EtcdSecKeyPrefix string
	EtcdProductKey   string
}

type RedisConfing struct {
	RedisAddr        string
	RedisPassword    string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type SecInfoConfing struct {
	ProductId int
	StartTime int
	EndTIme   int
	Status    int
	Total     int
	Left      int
}

func initConfig() (err error) {
	redisAddr := beego.AppConfig.String("redis_addr")
	redisPassword := beego.AppConfig.String("redis_password")

	redisMaxIdle, err := beego.AppConfig.Int("redis_max_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_idle error:%v", err)
		return
	}
	redisMaxActive, err := beego.AppConfig.Int("redis_max_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_max_active error:%v", err)
		return
	}
	redisIdleTimeout, err := beego.AppConfig.Int("redis_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_idle_timeout error:%v", err)
		return
	}

	etcdAddr := beego.AppConfig.String("etcd_addr")
	etcdSecKeyPrefix := beego.AppConfig.String("etcd_sec_key_prefix")
	etcdProductKey := beego.AppConfig.String("etcd_product_key")
	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")

	if err != nil {
		err = fmt.Errorf("init config failed, read EtcdTimeout error:%v", err)
		return
	}
	logs.Debug("read config successful,redisAddr:%v", redisAddr)
	logs.Debug("read config successful,etcdAddr:%v", etcdAddr)
	SeckillConf.RedisConf.RedisAddr = redisAddr
	SeckillConf.RedisConf.RedisMaxIdle = redisMaxIdle
	SeckillConf.RedisConf.RedisMaxActive = redisMaxActive
	SeckillConf.RedisConf.RedisIdleTimeout = redisIdleTimeout
	SeckillConf.RedisConf.RedisPassword = redisPassword

	SeckillConf.EtcdConf.EtcdAddr = etcdAddr
	SeckillConf.EtcdConf.EtcdTimeout = etcdTimeout
	SeckillConf.EtcdConf.EtcdSecKeyPrefix = etcdSecKeyPrefix

	if strings.HasSuffix(SeckillConf.EtcdConf.EtcdSecKeyPrefix, "/") == false {
		SeckillConf.EtcdConf.EtcdSecKeyPrefix = SeckillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}

	SeckillConf.EtcdConf.EtcdProductKey = fmt.Sprintf("%s%s", etcdSecKeyPrefix, etcdProductKey)

	if len(redisAddr) == 0 || len(etcdAddr) == 0 || len(redisPassword) == 0 || len(etcdSecKeyPrefix) == 0 || len(etcdProductKey) == 0 {
		err = fmt.Errorf("init config failed.redis[%s] or etcd[%s]  or redisPassword[%s] or etcdSecKeyPrefix[%s]  or etcdProductKey[%s] config is null", redisAddr, etcdAddr, redisPassword, etcdSecKeyPrefix, etcdProductKey)
		return
	}

	logpath := beego.AppConfig.String("log_path")
	loglevel := beego.AppConfig.String("log_level")
	SeckillConf.LogPath = logpath
	SeckillConf.LogLevel = loglevel

	return
}
