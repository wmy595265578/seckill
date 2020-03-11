package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"sync"
)

var (
	secLimitMgr = &SecLimitMgr{
		UserLimitMap: make(map[int]*Limit, 1000),
	}
	seclimit = &Limit{}
)

type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap   map[string]*Limit
	Lock         sync.Mutex
}

func antiSpam(req *SecRequest) (err error) {
	//secLimitMgr.Lock.Lock()
	//defer secLimitMgr.Lock.Unlock()
	//
	//secLimit, ok := secLimitMgr.UserLimitMap[req.UserId]
	//if !ok {
	//	secLimitMgr.UserLimitMap[req.UserId] = secLimit
	//}
	//if  seclimit.secLimit.Count(req.AccessTime.Unix()) > SeckillConf.UserSecAccessLimit {
	//	err = fmt.Errorf("invalid count ")
	//	return
	//}

	_, ok := SeckillConf.idBlackMap[req.UserId]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("userId[%v] is block by id black", req.UserId)
		return
	}

	_, ok = SeckillConf.ipBlackMap[req.ClientAddr]
	if ok {
		err = fmt.Errorf("invalid request")
		logs.Error("userId[%v] ip[%v] is block  by ip black", req.UserId, req.ClientAddr)
		return
	}

	SeckillConf.SecLimitMgrConf.Lock.Lock()
	//user Id
	limit, ok := SeckillConf.SecLimitMgrConf.UserLimitMap[req.UserId]
	if !ok {
		limit = &Limit{
			secLimit: &SecLimit{},
			minLimit: &MinLimit{},
		}
		SeckillConf.SecLimitMgrConf.UserLimitMap[req.UserId] = limit
	}

	secIdCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIdCount := limit.minLimit.Count(req.AccessTime.Unix())

	//user Ip

	limit, ok = SeckillConf.SecLimitMgrConf.IpLimitMap[req.ClientAddr]
	secIpCount := limit.secLimit.Count(req.AccessTime.Unix())
	minIpCount := limit.minLimit.Count(req.AccessTime.Unix())

	if secIdCount > SeckillConf.AccessLimitConf.UserSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if minIdCount > SeckillConf.AccessLimitConf.UserMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if secIpCount > SeckillConf.AccessLimitConf.IPSecAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	if minIpCount > SeckillConf.AccessLimitConf.IPMinAccessLimit {
		err = fmt.Errorf("invalid request")
		return
	}

	return
}
