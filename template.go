package main

import (
	"github.com/cheneylew/gocms/helper"
	"github.com/cheneylew/goutil/utils"
)

func TemplateMain()  {
	if false {
		text := helper.Pagination("http://www.baidu.com/",5,30,200)
		utils.JJKPrintln(text)
	}
}