package main

import (
	"github.com/astaxie/beego"
	_ "github.com/xzdbd/weixingis/routers"
)

func main() {
	beego.Info("come in")
	beego.Run()
}
