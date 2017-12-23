package controllers

import (
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/gocms/database"
	"github.com/cheneylew/gocms/models"
	"time"
	"github.com/cheneylew/gocms/conf"
)

type AdminCPController struct {
	BaseController
}

func (c *AdminCPController) Prepare() {
	c.BaseController.Prepare()
	c.Layout = "admin/layout.html"
	c.Data["Title"] = conf.GlobalConfig.ProductName+" | 控制台"

	menus := models.CreateMenus()
	activedMenuId, _ := c.GetInt64("active",0)
	if activedMenuId > 0 {
		for i := 0; i < len(menus); i++ {
			if menus[i].MenuID == activedMenuId {
				menus[i].IsActive = true
			} else {
				menus[i].IsActive = false
			}
		}
	}
	c.Data["Menus"] = menus

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
	c.TplName = "admin/dashboard.html"
}

func (c *AdminCPController) Login() {
	c.TplName = "admin/login.html"
	c.Data["IsLogin"] = true

	if c.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")
		user := database.DB.GetUserWithEmailOrUsername(username)
		md5Password := utils.MD5(password+user.UserSalt)
		if md5Password == user.Password {
			c.SaveUser(user)
			c.RedirectWithURLError("/admincp/home","登陆成功！")
		} else {
			c.SetError("用户名和密码不匹配", false)
		}
	}
}

func (c *AdminCPController) Regist() {
	c.TplName = "admin/regist.html"
	c.Data["IsLogin"] = true

	if c.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")

		user := database.DB.GetUserWithEmailOrUsername(username)
		if user != nil {
			c.SetError("用户已存在", false)
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
			c.SetError("注册用户失败", false)
			return
		}

		c.SetError("注册用户成功", true)
	}
}

func (c *AdminCPController) Logout() {
	c.UserLogout()
	c.RedirectWithURL("/admincp/login")
}