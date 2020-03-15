package models

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego/logs"
)

func userCheck(req *SecRequest) (err error) {
	found := false

	for _, refer := range SeckillConf.ReferWhiteList {
		if refer == req.ClientRefence {
			found = true
			break
		}
	}

	if !found {
		err = fmt.Errorf("invalid request")
		logs.Warn("user[%d] is reject by refer, req[%v]", req.UserId, req)
		return
	}

	authData := fmt.Sprintf("%d:%s",  req.UserId, SeckillConf.CookieSecretKey)
	authSin := fmt.Sprintf("%x", md5.Sum([]byte(authData)))

	if authSin != req.UserAuthSign {
		err = fmt.Errorf("invalid user cookie auth")
		return
	}

	return
}
