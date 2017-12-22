package main

import (
	_ "github.com/cheneylew/gocms/routers"
	"github.com/astaxie/beego"
)

func main() {
	MockMain()
	beego.Run()
}

