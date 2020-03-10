package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	RedisPool RedisPoolConf
)

type RedisPoolConf struct {
	RedisBlackConf       *redis.Pool
	RedisProxy2LayerConn *redis.Pool
	RedisLayer2ProxyConn *redis.Pool
}

func InitRedis() (err error) {
	err = initRedisBlackConf()
	if err != nil {
		logs.Error("initRedisBlackConf failed ,err:%v", err)
		return
	}

	err = initRedisLayer2ProxyConf()

	if err != nil {
		logs.Error("initRedisLayer2ProxyConf failed ,err:%v", err)
		return
	}

	err = initRedisProxy2LayerConf()
	if err != nil {
		logs.Error("initRedisProxy2LayerConf failed ,err:%v", err)
		return
	}
	return
}

func initRedisProxy2LayerConf() (err error) {
	redisPollConf := &redis.Pool{
		MaxIdle:     SeckillConf.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:   SeckillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(SeckillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", SeckillConf.RedisProxy2LayerConf.RedisAddr, redis.DialPassword(SeckillConf.RedisProxy2LayerConf.RedisPassword))
		},
	}
	conn := redisPollConf.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("conn redis failed,err:%v", err)
		return
	}
	RedisPool.RedisProxy2LayerConn = redisPollConf
	return
}

func initRedisBlackConf() (err error) {
	redisPollConf := &redis.Pool{
		MaxIdle:     SeckillConf.RedisBlackConf.RedisMaxIdle,
		MaxActive:   SeckillConf.RedisBlackConf.RedisMaxActive,
		IdleTimeout: time.Duration(SeckillConf.RedisBlackConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", SeckillConf.RedisBlackConf.RedisAddr, redis.DialPassword(SeckillConf.RedisBlackConf.RedisPassword))
		},
	}
	conn := redisPollConf.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("conn redis failed,err:%v", err)
		return
	}
	RedisPool.RedisBlackConf = redisPollConf
	return
}

func initRedisLayer2ProxyConf() (err error) {
	redisPollConf := &redis.Pool{
		MaxIdle:     SeckillConf.RedisProxy2LayerConf.RedisMaxIdle,
		MaxActive:   SeckillConf.RedisProxy2LayerConf.RedisMaxActive,
		IdleTimeout: time.Duration(SeckillConf.RedisProxy2LayerConf.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", SeckillConf.RedisProxy2LayerConf.RedisAddr, redis.DialPassword(SeckillConf.RedisProxy2LayerConf.RedisPassword))
		},
	}
	conn := redisPollConf.Get()
	defer conn.Close()
	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("conn redis failed,err:%v", err)
		return
	}
	RedisPool.RedisLayer2ProxyConn = redisPollConf
	return
}
