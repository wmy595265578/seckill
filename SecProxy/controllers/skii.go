package controllers

import "github.com/astaxie/beego"

type SkillController struct {
	beego.Controller
}


func (p *SkillController) SecKill() {
	p.Data["json"]="hello SecKill"
	p.ServeJSON()
}

func (p *SkillController)  SecInfo()  {
	p.Data["json"]="hello SecInfo"
	p.ServeJSON()
}
