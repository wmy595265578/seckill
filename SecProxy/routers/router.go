package routers

import (
	"seckill/SecProxy/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/seckill", &controllers.SkillController{},"*:SecKill")
    beego.Router("/secinfo", &controllers.SkillController{},"*:SecInfo")
}
