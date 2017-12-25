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
	c.AddCSS("dashboard.css")
	//c.AddJS("dashboard.js")

	c.Data["UserLogins"] = database.DB.GetUserLogins()
}

func (c *AdminCPController) Login() {
	c.AddCSS("login.css")
	c.TplName = "admin/login.html"
	c.Data["IsLogin"] = true

	utils.JJKPrintln(c.Ctx.Request.Header.Get("User-Agent"))

	if c.IsPost() {
		username := c.GetString("username")
		password := c.GetString("password")
		user := database.DB.GetUserWithEmailOrUsername(username)
		if user == nil {
			c.SetError("用户查询失败", false)
			return
		}
		md5Password := utils.MD5(password+user.UserSalt)
		if md5Password == user.Password {
			//login history
			login := &models.UserLogins{
				User:user,
				UserLoginDate:time.Now(),
				UserLoginIp:c.Ctx.Request.RemoteAddr,
				UserLoginBrowser:c.Ctx.Request.Header.Get("User-Agent"),

			}
			database.DB.Orm.Insert(login)

			//session
			c.SaveUser(user)
			c.RedirectWithURLError("/admincp/home","登陆成功！")
		} else {
			c.SetError("用户名和密码不匹配", false)
		}
	}
}

func (c *AdminCPController) Regist() {
	c.TplName = "admin/regist.html"
	c.AddCSS("login.css")
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

func (c *AdminCPController) Users() {
	isAddUser := c.Path(2) == "add"
	if isAddUser {
		c.TplName = "admin/users_add.html"
		c.Data["Roles"] = database.DB.GetUserRoles()

		if c.IsPost() {
			lastName := c.GetString("last_name","")
			username := c.GetString("username","")
			email := c.GetString("email","")
			password := c.GetString("password","")
			password2 := c.GetString("password2","")
			usergroups := c.GetString("usergroups[]","")
			firstName := c.GetString("first_name","")
			salt := utils.RandomString(32)

			u1 := database.DB.GetUserWithEmailOrUsername(username)
			u2 := database.DB.GetUserWithEmailOrUsername(email)
			if u1 != nil || u2 != nil {
				c.SetError("用户已经存在", false)
				return
			}

			if password != password2 {
				c.SetError("两次输入的密码不一致", false)
				return
			}

			user := &models.User{
				Username:username,
				UserEmail:email,
				Password:utils.MD5(password+salt),
				UserSalt:salt,
				UserFirstName:firstName,
				UserLastName:lastName,
				UserRegistDate:time.Now(),
				UserRole:&models.UserRole{
					UserRoleId:utils.JKStrToInt64(usergroups),
				},
			}

			a, e := database.DB.Orm.Insert(user)
			if e != nil || a == 0 {
				c.SetError("注册失败", false)
			} else {
				c.RedirectWithURLError("/admincp/users?active=3","注册成功！")
			}

		}
	} else {
		c.TplName = "admin/users_manage.html"
		c.AddCSS("dataset.css")

		c.Data["Users"] = database.DB.GetUsers()
	}
}

func (c *AdminCPController) Types() {
	fun := c.Path(2)
	if fun == "new" {
		c.TplName = "admin/types_new.html"
		if c.IsPost() {
			isPrivileged, _ := c.GetBool("is_privileged", false)
			template := c.GetString("template","")
			baseUrl := c.GetString("base_url","")
			name := c.GetString("name","")
			isStandard, _ := c.GetBool("is_standard", false)

			systemName := utils.SnakeString(name)
			//创建表
			database.DB.DBBaseCreateTableWithContentID(systemName)

			//记录content type
			contentType := &models.ContentType{
				SystemName:systemName,
				Name:utils.UpperFirstChar(name),
				IsPrivileged:isPrivileged,
				IsStandard:isStandard,
				BaseUrl:baseUrl,
				Template:template,
			}

			a, e := database.DB.Orm.Insert(contentType)
			if e != nil || a == 0 {
				c.SetError("创建失败", false)
			} else {
				c.SetError("创建成功", true)
			}
		}
	} else if fun == "manage" {
		c.TplName = "admin/types_manage.html"
		c.AddCSS("dataset.css")

		contentType := database.DB.GetContentTypeWithId(c.PathInt64(3))
		c.Data["ContentType"] = contentType
	} else if fun == "edit" {
		c.TplName = "admin/types_manage.html"
		c.AddCSS("dataset.css")
	} else {
		c.TplName = "admin/types.html"
		c.AddCSS("dataset.css")
		c.Data["ContentTypes"] = database.DB.GetContentTypes()
	}
}

func (c *AdminCPController) Fields() {
	c.TplName = "admin/fields_add.html"
	c.AddJS("form.field.js")
	c.Data["ContentType"] = database.DB.GetContentTypeWithId(c.PathInt64(2))
}

func (c *AdminCPController) Custom_fields() {
	fun := c.Path(2)
	if fun == "ajax_field_form" {
		ftype := c.GetString("type","")
		html := models.GetFieldTypeHTML(ftype)
		if ftype == "checkbox" {
			html := &models.FieldTypeCheckBox{
				FieldType:models.FieldType{},
			}
			c.Ctx.WriteString(html.ToHTML())
		} else {
			c.Ctx.WriteString(html)
		}
	}
}