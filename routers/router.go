package routers

import (
	"github.com/astaxie/beego"
	"github.com/xzdbd/weixingis/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/weixin", &controllers.MainController{})
}
