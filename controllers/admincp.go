package controllers

import (
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/gocms/database"
	"github.com/cheneylew/gocms/models"
	"time"
)

type AdminCPController struct {
	BaseController
}

func (c *AdminCPController) Prepare() {
	c.BaseController.Prepare()
	c.Layout = "admin/layout.html"

	//登陆检测
	urlPath := c.Ctx.Request.URL.Path

	var whitePaths []string
	whitePaths = append(whitePaths, "/admincp/login")
	whitePaths = append(whitePaths, "/admincp/regist")

	if !utils.Contain(whitePaths, urlPath) {
		if c.IsLogin() {
			c.Data["User"] = c.GetUser()
		} else {
			c.RedirectWithURL("/admincp/login")
		}
	}
}

func (c *AdminCPController) Finish() {
	c.Controller.Finish()
}

func (c *AdminCPController) Index() {
	c.TplName = "main.html"
}


func (c *AdminCPController) Home() {
	c.TplName = "admin/home.html"
}

func (c *AdminCPController) Login() {
	c.TplName = "admin/login.html"

	if c.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")
		user := database.DB.GetUserWithEmailOrUsername(username)
		md5Password := utils.MD5(password+user.UserSalt)
		if md5Password == user.Password {
			c.SaveUser(user)
			c.RedirectWithURLError("/admincp/home","登陆成功！")
		} else {
			c.SetError("用户名和密码不匹配")
		}
	}
}

func (c *AdminCPController) Regist() {
	c.TplName = "admin/regist.html"
	if c.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")

		user := database.DB.GetUserWithEmailOrUsername(username)
		if user != nil {
			c.SetError("用户已存在")
			return ;
		}

		salt := utils.RandomString(32)
		md5Pwd := utils.MD5(password+salt)

		role := database.DB.GetUserRolesWithGrade(0)
		newUser := &models.User{
			UserRole:role,
			Username:username,
			Password:md5Pwd,
			UserRegistDate:time.Now(),
			UserSalt:salt,
		}

		a, e := database.DB.Orm.Insert(newUser)
		if e != nil || a == 0 {
			utils.JJKPrintln(e)
			c.SetError("注册用户失败")
			return
		}

		c.SetError("注册用户成功")
	}
}