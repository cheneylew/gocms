package main

import (
	"github.com/cheneylew/gocms/helper"
	"github.com/cheneylew/goutil/utils"
)

func TemplateMain()  {
	value := `男=1
 女=2
 中性=3
 无性别=4`
	options := helper.StringToOptions(value)
	utils.JJKPrintln(utils.JKJSON(options))

	utils.JJKPrintln(helper.OptionsToString(options))

}