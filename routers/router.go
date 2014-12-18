package routers

import (
	"github.com/xzdbd/weixingis/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
