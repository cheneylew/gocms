package controllers

import (
	"github.com/cheneylew/goutil/utils"
	beego2 "github.com/cheneylew/goutil/utils/beego"
	"github.com/cheneylew/gocms/conf"
	"fmt"
)

var BEEGO_CONFIG beego2.BeegoConfig
var FILTER_PATHS  []string

func init() {
	// configs
	FILTER_PATHS = append(FILTER_PATHS,"/user/login")
	FILTER_PATHS = append(FILTER_PATHS,"/user/regist")
	BEEGO_CONFIG = beego2.BeegoConfig{
		LoginCheck:false,
	}

}

type BaseController struct {
	beego2.BBaseController
}

func (c *BaseController) Prepare() {
	c.BBaseController.Prepare()
	urlPath := c.Ctx.Request.URL.Path
	c.Data["Config"] = conf.GlobalConfig

	error := c.GetString("msg")
	if len(error) > 0 {
		c.SetError(error, true)
	}

	c.Layout = "layout.html"
	if BEEGO_CONFIG.LoginCheck {
		if !utils.Contain(FILTER_PATHS, urlPath) {
			if c.IsLogin() {
				c.Data["User"] = c.GetUser()
			} else {
				c.RedirectWithURL("/user/login")
			}
		}
	}
}

func (c *BaseController) SetError(error string, isNotice bool) {
	if isNotice {
		c.Data["Error"] = fmt.Sprintf("<div class=\"notice\">%s</div>",error)
	} else {
		c.Data["Error"] = fmt.Sprintf("<div class=\"error\">%s</div>",error)
	}
}

func (c *BaseController) RedirectWithURLError(url, error string) {
	c.RedirectWithURL(url+"?msg="+error)
}




