package routers

import (
	"gocms.win/cheneylew/gocms/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
}
