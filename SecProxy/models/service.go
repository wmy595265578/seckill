package models

import (
	"fmt"
	"seckill/SecProxy/conf"
)

func SecInfo(productId int) (data map[string]interface{}, code int, err error) {
	data = make(map[string]interface{}, 16)
	conf.SeckillConf.RWSecKillLock.RLock()
	defer conf.SeckillConf.RWSecKillLock.RUnlock()
	v, ok := conf.SeckillConf.SecInfoConfMap[productId]
	if !ok {
		code = ErrInvalidRequest
		err = fmt.Errorf("not found product_id:%d", productId)
		return
	}
	data["product_id"] = v.ProductId
	data["start_time"] = v.StartTime
	data["end_time"] = v.EndTIme
	data["status"] = v.Status
	return
}
