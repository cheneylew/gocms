package main

import (
	_ "github.com/cheneylew/gocms/routers"
	"github.com/astaxie/beego"
	"github.com/cheneylew/goutil/utils"
)

func main() {
	MockMain()
	TemplateMain()
	beego.AddFuncMap("Equal",utils.Equal)
	beego.AddFuncMap("InSlice",utils.InSlice)
	beego.AddFuncMap("ToStr",utils.ToString)
	beego.Run()
}

