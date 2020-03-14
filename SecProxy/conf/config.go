package conf

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"seckill/SecProxy/models"
	"strings"
)

func initConfig() (err error) {
	//redis 接入层->业务逻辑层
	redisProxy2layerAddr := beego.AppConfig.String("redis_proxy2layer_addr")
	redisProxy2layerPassword := beego.AppConfig.String("redis_proxy2layer_password")

	redisProxy2layerMaxIdle, err := beego.AppConfig.Int("redis_proxy2layer_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle error:%v", err)
		return
	}

	redisProxy2layerMaxActive, err := beego.AppConfig.Int("redis_proxy2layer_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_active error:%v", err)
		return
	}
	redisProxy2layerIdleTimeout, err := beego.AppConfig.Int("redis_proxy2layer_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_proxy2layer_idle_timeout error:%v", err)
		return
	}

	models.SeckillConf.RedisProxy2LayerConf.RedisAddr = redisProxy2layerAddr
	models.SeckillConf.RedisProxy2LayerConf.RedisPassword = redisProxy2layerPassword
	models.SeckillConf.RedisProxy2LayerConf.RedisMaxActive = redisProxy2layerMaxActive
	models.SeckillConf.RedisProxy2LayerConf.RedisIdleTimeout = redisProxy2layerIdleTimeout
	models.SeckillConf.RedisProxy2LayerConf.RedisMaxIdle = redisProxy2layerMaxIdle

	//redis黑名单相关配置
	redisBlackAddr := beego.AppConfig.String("redis_black_addr")
	redisBlackPassword := beego.AppConfig.String("redis_black_password")

	redisBlackMaxIdle, err := beego.AppConfig.Int("redis_black_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_idle error:%v", err)
		return
	}
	redisBlackMaxActive, err := beego.AppConfig.Int("redis_black_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_active error:%v", err)
		return
	}
	redisBlackIdleTimeout, err := beego.AppConfig.Int("redis_black_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_black_idle_timeout error:%v", err)
		return
	}

	models.SeckillConf.RedisBlackConf.RedisAddr = redisBlackAddr
	models.SeckillConf.RedisBlackConf.RedisPassword = redisBlackPassword
	models.SeckillConf.RedisBlackConf.RedisMaxActive = redisBlackMaxActive
	models.SeckillConf.RedisBlackConf.RedisIdleTimeout = redisBlackIdleTimeout
	models.SeckillConf.RedisBlackConf.RedisMaxIdle = redisBlackMaxIdle

	//redis 业务逻辑层->接入层
	redisLayer2proxyAddr := beego.AppConfig.String("redis_layer2proxy_addr")
	redisLayer2proxyPassword := beego.AppConfig.String("redis_layer2proxy_password")

	redisLayer2proxyMaxIdle, err := beego.AppConfig.Int("redis_layer2proxy_idle")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_idle error:%v", err)
		return
	}
	redisLayer2proxyMaxActive, err := beego.AppConfig.Int("redis_layer2proxy_active")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_active error:%v", err)
		return
	}
	redisLayer2proxyIdleTimeout, err := beego.AppConfig.Int("redis_layer2proxy_idle_timeout")
	if err != nil {
		err = fmt.Errorf("init config failed, read redis_layer2proxy_idle_timeout error:%v", err)
		return
	}

	models.SeckillConf.RedisLayer2ProxyConf.RedisAddr = redisLayer2proxyAddr
	models.SeckillConf.RedisLayer2ProxyConf.RedisPassword = redisLayer2proxyPassword
	models.SeckillConf.RedisLayer2ProxyConf.RedisMaxActive = redisLayer2proxyMaxActive
	models.SeckillConf.RedisLayer2ProxyConf.RedisIdleTimeout = redisLayer2proxyIdleTimeout
	models.SeckillConf.RedisLayer2ProxyConf.RedisMaxIdle = redisLayer2proxyMaxIdle

	etcdAddr := beego.AppConfig.String("etcd_addr")
	etcdSecKeyPrefix := beego.AppConfig.String("etcd_sec_key_prefix")
	etcdProductKey := beego.AppConfig.String("etcd_product_key")
	etcdTimeout, err := beego.AppConfig.Int("etcd_timeout")

	if err != nil {
		err = fmt.Errorf("init config failed, read EtcdTimeout error:%v", err)
		return
	}
	//logs.Debug("read config successful,redisAddr:%v", redisAddr)
	logs.Debug("read config successful,etcdAddr:%v", etcdAddr)

	models.SeckillConf.EtcdConf.EtcdAddr = etcdAddr
	models.SeckillConf.EtcdConf.EtcdTimeout = etcdTimeout
	models.SeckillConf.EtcdConf.EtcdSecKeyPrefix = etcdSecKeyPrefix

	if strings.HasSuffix(models.SeckillConf.EtcdConf.EtcdSecKeyPrefix, "/") == false {
		models.SeckillConf.EtcdConf.EtcdSecKeyPrefix = models.SeckillConf.EtcdConf.EtcdSecKeyPrefix + "/"
	}

	models.SeckillConf.EtcdConf.EtcdProductKey = fmt.Sprintf("%s%s", models.SeckillConf.EtcdConf.EtcdSecKeyPrefix, etcdProductKey)

	//if len(redisAddr) == 0 || len(etcdAddr) == 0 || len(redisPassword) == 0 || len(etcdSecKeyPrefix) == 0 || len(etcdProductKey) == 0 {
	//	err = fmt.Errorf("init config failed.redis[%s] or etcd[%s]  or redisPassword[%s] or etcdSecKeyPrefix[%s]  or etcdProductKey[%s] config is null", redisAddr, etcdAddr, redisPassword, etcdSecKeyPrefix, etcdProductKey)
	//	return
	//}

	logpath := beego.AppConfig.String("log_path")
	loglevel := beego.AppConfig.String("log_level")
	models.SeckillConf.LogPath = logpath
	models.SeckillConf.LogLevel = loglevel

	models.SeckillConf.CookieSecretKey = beego.AppConfig.String("cookie_secretkey")

	referList := beego.AppConfig.String("refer_whitelist")
	if len(referList) > 0 {
		models.SeckillConf.ReferWhiteList = strings.Split(referList, ",")
	}

	ipSecLimit, err := beego.AppConfig.Int("ip_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_sec_access_limit error:%v", err)
		return
	}
	models.SeckillConf.AccessLimitConf.IPSecAccessLimit = ipSecLimit

	UserSecAccessLimit, err := beego.AppConfig.Int("user_sec_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_sec_access_limit error:%v", err)
		return
	}
	models.SeckillConf.UserSecAccessLimit = UserSecAccessLimit

	UserIdMinLimit, err := beego.AppConfig.Int("user_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read user_min_access_limit error:%v", err)
		return
	}

	models.SeckillConf.AccessLimitConf.UserMinAccessLimit = UserIdMinLimit
	ipMinLimit, err := beego.AppConfig.Int("ip_min_access_limit")
	if err != nil {
		err = fmt.Errorf("init config failed, read ip_min_access_limit error:%v", err)
		return
	}

	models.SeckillConf.AccessLimitConf.IPMinAccessLimit = ipMinLimit
	return
}
