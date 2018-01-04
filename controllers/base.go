package controllers

import (
	"github.com/cheneylew/goutil/utils"
	beego2 "github.com/cheneylew/goutil/utils/beego"
	"github.com/cheneylew/gocms/conf"
	"fmt"
	"strings"
	"time"
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

	forms,_ := utils.ToJSON(c.PostForm())
	if len(forms) > 0 && forms != "null" {
		utils.JJKPrintln(forms)
		codes := ""
		for key, _ := range c.PostForm() {
			varKey := utils.LowerFirstChar(utils.CamelString(key))
			if strings.Contains(key,"is") {
				codes += fmt.Sprintf("%s, _ := c.GetBool(\"%s\", false)\n", varKey, key)
			} else {
				codes += fmt.Sprintf("%s := c.GetString(\"%s\",\"\")\n", varKey, key)
			}

		}
		utils.JJKPrintln(codes)
	}

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
	decodeStr := utils.Base64Decode(error)
	if len(decodeStr) > 0 {
		error = decodeStr
	}
	if isNotice {
		c.Data["Error"] = fmt.Sprintf("<div class=\"notice\">%s</div>",error)
	} else {
		c.Data["Error"] = fmt.Sprintf("<div class=\"error\">%s</div>",error)
	}
}

func (c *BaseController) RedirectWithURLError(url, error string) {
	c.RedirectWithURL(url+"?msg="+utils.Base64Encode(error))
}

func (c *BaseController) AddCSS(path string) {
	if !strings.HasPrefix(path, "/") {
		path = "/static/branding/default/css/"+path
	}
	if c.Data["CSS"] == nil {
		c.Data["CSS"] = []string{path}
	} else {
		c.Data["CSS"] = append(c.Data["CSS"].([]string), path)
	}
}

func (c *BaseController) AddJS(path string) {
	if !strings.HasPrefix(path, "/") {
		path = "/static/branding/default/js/"+path
	}
	if c.Data["JS"] == nil {
		c.Data["JS"] = []string{path}
	} else {
		c.Data["JS"] = append(c.Data["JS"].([]string), path)
	}
}

func (c *BaseController) GetDateTime(key string) time.Time {
	date := c.GetString(key,"")
	dateHour := c.GetString(key+"_hour","")
	dateMinute := c.GetString(key+"_minute","")
	dateAmpm := c.GetString(key+"_ampm","")

	datetime := utils.ValuesToDateTime(date, dateHour, dateMinute, dateAmpm)
	return datetime
}

func (c *BaseController) GetDateTimeStr(key string) string {
	datetime := c.GetDateTime(key)
	return utils.JKTimeToString(datetime)
}

func (c *BaseController) GetLimit() int64 {
	limit , _ := c.GetInt64("limit", conf.GlobalConfig.PageLimit)
	return limit
}

func (c *BaseController) GetOffset() int64 {
	offset , _ := c.GetInt64("offset", 0)
	return offset
}





