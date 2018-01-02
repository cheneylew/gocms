package main

import (
	"fmt"
	"github.com/cheneylew/goutil/utils"
)

func TemplateMain()  {

	var hours []string
	var minutes []string
	for i := 0; i<= 12; i++ {
		hours = append(hours, fmt.Sprintf("%02d", i))
	}
	for i := 0; i<= 59; i++ {
		minutes = append(minutes, fmt.Sprint("%02d", i))
	}

	utils.JJKPrintln(hours)
}
