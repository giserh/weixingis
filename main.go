package main

import (
	"github.com/astaxie/beego"
	_ "github.com/xzdbd/weixingis/routers"
)

func main() {
	beego.SetLogger("smtp", `{"username":"weixingis@163.com","password":"gjfeixphowfkqhlb","host":"smtp.163.com:25","fromAddress":"weixingis@163.com","subject":"Log from weixingis","sendTos":["xzdbd@sina.com"]}`)
	beego.Info("come in")
	beego.Run()
}
