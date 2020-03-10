package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

func SecInfoList() (data []map[string]interface{}, code int, err error) {
	SeckillConf.RWSecKillLock.RLock()
	defer SeckillConf.RWSecKillLock.RUnlock()
	for _, v := range SeckillConf.SecInfoConfMap {
		item, _, err := SecInfoById(v.ProductId)
		if err != nil {
			logs.Error("get product_id[%d] failed,err:%v", err)
			continue
		}
		logs.Debug("get product[%d],result[%v],all[%v],v[%v]", v.ProductId, item, SeckillConf.SecInfoConfMap, v)
		data = append(data, item)
	}
	return
}

func SecInfo(productId int) (data []map[string]interface{}, code int, err error) {

	SeckillConf.RWSecKillLock.RLock()
	defer SeckillConf.RWSecKillLock.RUnlock()
	item, code, err := SecInfoById(productId)
	if err != nil {
		return
	}
	data = append(data, item)
	return
}

func SecInfoById(productId int) (data map[string]interface{}, code int, err error) {
	SeckillConf.RWSecKillLock.RLock()
	defer SeckillConf.RWSecKillLock.RUnlock()

	v, ok := SeckillConf.SecInfoConfMap[productId]
	if !ok {
		code = ErrNotFoundProductId
		err = fmt.Errorf("not found produt_id:%d", productId)
		return
	}
	start := false
	end := false
	status := "success"

	now := time.Now().Unix()

	if now-v.StartTime > 0 {
		start = true
	}
	if now-v.StartTime < 0 {
		start = false
		end = false
		status = "sec kill is not start"
		code = ErrActiveNotStart
	}
	if now-v.EndTIme > 0 {
		start = false
		end = true
		status = "sec kill is already end"
		code = ErrActiveAlreadyEnd
	}

	if v.Status == ProductStatusForceSaleOut || v.Status == ProductStatusSaleOut {
		start = false
		end = true
		status = "product is sale out"
		code = ErrActiveSaleOut
	}

	data = make(map[string]interface{}, 16)
	data["product_id"] = v.ProductId
	data["start"] = start
	data["end"] = end
	data["status"] = status
	return
}

func SecKill(req *SecRequest) (data map[string]interface{}, code int, err error) {
	SeckillConf.RWSecKillLock.RLock()
	defer SeckillConf.RWSecKillLock.RUnlock()

	err = userCheck(req)
	if err != nil {
		code = ErrUserCheckAuthFailed
		logs.Warn("userId[%d], invalid,checked failed,req[%v]", req.UserId, req)
		return
	}

	err = antiSpam(req)
	if err != nil {
		code = ErrUserServiceBusy
		logs.Warn("userId[%d] invalid, check failed, req[%v]", req.UserId, req)
		return
	}
	data, code, err = SecInfoById(req.ProductId)
	if err != nil {
		logs.Warn("userId[%d] secInfoBy Id failed, req[%v]", req.UserId, req)
		return
	}

	if code != 0 {
		logs.Warn("userId[%d] secInfoByid failed, code[%d] req[%v]", req.UserId, code, req)
		return
	}

	userKey := fmt.Sprintf("%s %s", req.UserId, req.ProductId)
	SeckillConf.UserConnMap[userKey] = req.ResultChan

	SeckillConf.SecReqChan <- req

	ticker := time.NewTicker(time.Second * 10)

	defer func() {
		ticker.Stop()
		SeckillConf.UserConnMapLock.Lock()
		delete(SeckillConf.UserConnMap, userKey)
		SeckillConf.UserConnMapLock.Unlock()
	}()

	select {
	case <-ticker.C:
		code = ErrProcessTimeout
		err = fmt.Errorf("request timeout")
		return
	case <-req.CloseNotify:
		code = ErrClientClosed
		err = fmt.Errorf("client already closed")

	case result := <-req.ResultChan:
		code = result.Code
		data["product_id"] = result.ProductId
		data["token"] = result.Token
		data["user_id"] = result.UserId
		return
	}

	return
}
