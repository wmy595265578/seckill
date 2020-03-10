package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

func InitService(serviceConf *SecSkillConf) (err error) {
	err = loadBlackList()
	if err != nil {
		logs.Error("load black list err:%v", err)
		return
	}

	logs.Debug("init service succ, config:%v", SeckillConf)

	SeckillConf.SecLimitMgrConf = &SecLimitMgr{
		UserLimitMap: make(map[int]*Limit, 10000),
		IpLimitMap:   make(map[string]*Limit, 10000),
	}

	SeckillConf.SecReqChan = make(chan *SecRequest, SeckillConf.SecReqChanSize)
	SeckillConf.UserConnMap = make(map[string]chan *SecResult, 100000)

	initRedisProcessFunc()

	return
}

func initRedisProcessFunc() {
	for i := 0; i < SeckillConf.WriteProxy2LayerGoroutineNum; i++ {
		go WriteHandle()
	}

	for i := 0; i < SeckillConf.ReadProxy2LayerGoroutineNum; i++ {
		go ReadHandle()
	}
}

func loadBlackList() (err error) {
	SeckillConf.ipBlackMap = make(map[string]bool, 10000)
	SeckillConf.idBlackMap = make(map[int]bool, 10000)

	conn := RedisPool.RedisBlackConf.Get()
	defer conn.Close()

	reply, err := conn.Do("hgetall", "idblacklist")

	idlist, err := redis.Strings(reply, err)

	if err != nil {
		logs.Warn("hgetall idblacklist failed,err:%v", err)
		return
	}
	for _, v := range idlist {
		id, err := strconv.Atoi(v)
		if err != nil {
			logs.Warn("invalid user id [%v],", err)
			continue
		}
		SeckillConf.idBlackMap[id] = true
	}

	reply, err = conn.Do("hgetall", "ipblacklist")
	iplist, err := redis.Strings(reply, err)

	if err != nil {
		logs.Warn("hgetall ipblacklist failed,err:%v", err)
		return
	}
	for _, v := range iplist {
		SeckillConf.ipBlackMap[v] = true
	}

	go SyncIpBlackList()
	go SyncIdBlackList()
	return
}

func SyncIpBlackList() {
	var ipList []string
	lastTime := time.Now().Unix()

	for {
		conn := RedisPool.RedisBlackConf.Get()
		defer conn.Close()

		reply, err := conn.Do("BLPOP", "blackiplist", time.Second)
		ip, err := redis.String(reply, err)

		if err != nil {
			continue
		}
		ipList = append(ipList, ip)

		curTime := time.Now().Unix()
		if len(ipList) > 100 || curTime-lastTime > 5 {
			SeckillConf.RWBlackLock.Lock()
			for _, v := range ipList {
				SeckillConf.ipBlackMap[v] = true
			}
			SeckillConf.RWBlackLock.Unlock()
			lastTime = curTime
			logs.Info("sync ip list from redis succ, ip[%v]", ipList)

		}
	}
}

func SyncIdBlackList() {
	for {
		conn := RedisPool.RedisBlackConf.Get()
		defer conn.Close()

		reply, err := conn.Do("BLPOP", "blackidlist", time.Second)
		id, err := redis.Int(reply, err)
		if err != nil {
			continue
		}
		SeckillConf.RWBlackLock.Lock()
		SeckillConf.idBlackMap[id] = true
		SeckillConf.RWBlackLock.Unlock()
		logs.Info("sync id list from redis succ, ip[%v]", id)

	}
}
