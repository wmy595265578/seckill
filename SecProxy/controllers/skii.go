package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"seckill/SecProxy/models"
)

type SkillController struct {
	beego.Controller
}

func (p *SkillController) SecKill() {
	p.Data["json"] = "hello SecKill"
	p.ServeJSON()
}

func (p *SkillController) SecInfo() {
	productId, err := p.GetInt("product_id")
	result := make(map[string]interface{})
	result["code"] = 0
	result["message"] = "successful"
	if err != nil {

		result["code"] = 1001
		result["message"] = "invalid product_id"
		p.Data["json"] = result
		p.ServeJSON()
		logs.Error("invalid request,get product_id failed ,err:%v", err)
		return
	}

	data, code, err := models.SecInfo(productId)
	if err != nil {
		result["code"] = code
		result["message"] = err.Error()
		p.Data["json"] = result
		p.ServeJSON()
		logs.Error("invalid request,return product_id failed ,err:%v", err)

		return
	}
	p.Data["json"] = data
	p.ServeJSON()
}
