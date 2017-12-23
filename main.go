package main

import (
	_ "github.com/cheneylew/gocms/routers"
	"github.com/astaxie/beego"
	"github.com/cheneylew/goutil/utils"
)

func main() {
	MockMain()
	beego.AddFuncMap("Equal",utils.Equal)
	beego.Run()
}

