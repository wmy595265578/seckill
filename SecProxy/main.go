package main

import (
	"github.com/astaxie/beego"
	"seckill/SecProxy/conf"
	_ "seckill/SecProxy/routers"
)

func main() {
	err := conf.InitSec()
	if err != nil {
		panic(err)
		return
	}
	beego.Run()
}
