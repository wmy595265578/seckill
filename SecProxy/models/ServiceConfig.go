package models

import (
	"sync"
	"time"
)

const (
	ProductStatusNormal       = 0
	ProductStatusSaleOut      = 1
	ProductStatusForceSaleOut = 2
)

var (
	SeckillConf = &SecSkillConf{
		SecInfoConfMap: make(map[int]*SecInfoConfing),
	}
)

type SecSkillConf struct {
	RedisBlackConf       RedisConfing
	RedisProxy2LayerConf RedisConfing
	RedisLayer2ProxyConf RedisConfing
	EtcdConf             EtcdConfig
	LogPath              string
	LogLevel             string
	SecInfoConfMap       map[int]*SecInfoConfing
	RWSecKillLock        sync.RWMutex
	UserSecAccessLimit   int
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
	StartTime int64
	EndTIme   int64
	Status    int
	Total     int
	Left      int
}

type SecRequest struct {
	ProductId    int
	Source       string
	AuthCode     string
	SecTime      string
	Nance        string
	UserId       int
	UserAuthSign string
	AccessTime   time.Time
}
