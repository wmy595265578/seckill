package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"seckill/SecProxy/models"
	"strconv"
	"time"
	"fmt"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})

	result["code"] = 0
	result["message"] = "successful"
	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()

	if err != nil {
		result["code"] = models.ErrInvalidRequest
		result["message"] = "invalid product_id"
		return
	}
	source := p.GetString("source")
	authcode := p.GetString("authcode")
	secTime := p.GetString("time")
	nance := p.GetString("nance")

	secRequest := &models.SecRequest{}
	secRequest.ProductId = productId
	secRequest.AuthCode = authcode
	secRequest.Nance = nance
	secRequest.SecTime = secTime
	secRequest.Source = source
	secRequest.UserId, err = strconv.Atoi(p.Ctx.GetCookie("userId"))
	secRequest.UserAuthSign = p.Ctx.GetCookie("UserAuthSign")
	secRequest.AccessTime = time.Now()
	logs.Debug("client request:[%v]", secRequest)

	if err != nil {
		result["code"] = models.ErrInvalidRequest
		result["message"] = fmt.Sprintf("invalid cookie:userId")
		return
	}
	data, code, err := models.SecKill(secRequest)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		return
	}
	result["data"] = data
	result["code"] = code

}

func (p *SkillController) SecInfo() {
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "successful"

	defer func() {
		p.Data["json"] = result
		p.ServeJSON()
	}()
	if err != nil {
		data, code, err := models.SecInfoList()
		if err != nil {
			result["code"] = code
			result["message"] = err.Error()
			logs.Error("invalid request,get product_id failed ,err:%v", err)
			return
		}
		result["code"] = code
		result["data"] = data

	} else {
		data, code, err := models.SecInfo(productId)
		if err != nil {
			result["code"] = code
			result["message"] = err.Error()
			logs.Error("invalid request,return product_id failed ,err:%v", err)
			return
		}
		result["data"] = data
		result["code"] = code
		logs.Debug("get data from product_id ,data[%v]", data)
	}
}
