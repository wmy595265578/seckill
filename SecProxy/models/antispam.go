package models

import (
	"fmt"
	"sync"
)

var (
	secLimitMgr = &SecLimitMgr{
		UserLimitMap:make(map[int]*Limit,1000),
	}
	seclimit = &Limit{}
)

type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	Lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {
	secLimitMgr.Lock.Lock()
	defer secLimitMgr.Lock.Unlock()

	secLimit, ok := secLimitMgr.UserLimitMap[req.UserId]
	if !ok {
		secLimitMgr.UserLimitMap[req.UserId] = secLimit
	}
	if  seclimit.secLimit.Count(req.AccessTime.Unix()) > SeckillConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid count ")
		return
	}
	return
}
