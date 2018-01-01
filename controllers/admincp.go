package controllers

import (
	"github.com/cheneylew/goutil/utils"
	"github.com/cheneylew/gocms/database"
	"github.com/cheneylew/gocms/models"
	"time"
	"github.com/cheneylew/gocms/conf"
	"strings"
	"github.com/cheneylew/gocms/helper"
	"fmt"
)

type AdminCPController struct {
	BaseController
}

func (c *AdminCPController) Prepare() {
	c.BaseController.Prepare()
	c.Layout = "admin/layout.html"
	c.Data["Title"] = conf.GlobalConfig.ProductName+" | 控制台"

	menus := helper.CreateMenus()
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
			c.RedirectWithURL(fmt.Sprintf("/admincp/login?redirect_url=%s", utils.Base64Encode(c.Ctx.Request.RequestURI)))
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

	redirect_url := c.GetString("redirect_url","")
	c.Data["redirect_url"] = redirect_url

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

			if len(redirect_url) > 0 {
				c.RedirectWithURL(utils.Base64Decode(redirect_url))
			} else {
				c.RedirectWithURLError("/admincp/home","登陆成功！")
			}
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
		fields := database.DB.GetFieldTypesWithContentTypeId(contentType.ContentTypeId)
		c.Data["Fields"] = fields
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
	tableInfo := database.DB.GetContentTypeWithId(c.PathInt64(2))
	c.TplName = "admin/fields_add.html"
	c.AddJS("form.field.js")
	c.Data["ContentType"] = tableInfo

	if c.IsPost() {
		isRequired, _ := c.GetBool("required",false)
		fieldName := c.GetString("name","")
		helpText := c.GetString("help","")
		fieldType := c.GetString("type","")
		fieldTypeDefaultValue := c.GetString("default","")

		systemName := utils.SnakeString(fieldName)

		a := database.DB.GetFieldTypesWithContentTypeId(tableInfo.ContentTypeId)

		isExist := false

		for _, value := range a {
			if value.Name == systemName {
				isExist = true
			}
		}

		if isExist {
			c.SetError("该字段名称已存在", false)
			return
		}

		model := &models.FieldType{
			ContentType:tableInfo,
			SystemName:systemName,
			Name:fieldName,
			Help:helpText,
			Type:fieldType,
			Required:isRequired,
			DefaultValue:fieldTypeDefaultValue,
		}

		database.DB.DBBaseAddColumnTinyInt(tableInfo.SystemName, systemName)
		_, err := database.DB.Orm.Insert(model)
		if err != nil {
			utils.JJKPrintln(err)
		} else {
			utils.JJKPrintln("field add ok!")
		}
}
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
		} else if ftype == "date" {
			html := &models.FieldTypeDate{
				FieldType:models.FieldType{},
			}
			c.Ctx.WriteString(html.ToHTML())
		}  else if ftype == "datetime" {
			html := &models.FieldTypeDateTime{
				FieldType:models.FieldType{},
			}
			c.Ctx.WriteString(html.ToHTML())
		} else {
			c.Ctx.WriteString(html)
		}
	}
}

func (c *AdminCPController) Dataset() {
	fun := c.Path(2)
	if fun == "prep_actions" {
		//删除字段请求
		items := c.GetString("items","")
		var itemIds []int64
		for _, value := range strings.Split(items, "&") {
			kv := strings.Split(value, "=")
			if len(kv) == 2 {
				itemIds = append(itemIds, utils.JKStrToInt64(kv[1]))
			}
		}
		//returnUrl := c.GetString("return_url","")
		action := c.GetString("action");
		for _, value := range itemIds {
			if action == "delete_articles" {
				//删除文章
				item := database.DB.GetContentWithContentID(value)
				database.DB.Orm.Delete(item)
			} else if action == "delete_fieldTypes" {
				//删除自定义字段
				fieldType := database.DB.GetFieldTypesWithFieldTypeID(value)
				database.DB.Orm.Delete(fieldType)
				database.DB.DBBaseDropColumn(fieldType.ContentType.SystemName, fieldType.SystemName)
			} else if action == "delete_content_types" {
				contentType := database.DB.GetContentTypeWithId(value)

				//删除contentType
				database.DB.Orm.Delete(contentType)

				//删除content
				types := database.DB.GetContentsWithContentTypeID(contentType.ContentTypeId)
				for _, type1 := range types {
					database.DB.Orm.Delete(type1)
				}

				//删除自定义表
				database.DB.DBBaseDropTable(contentType.SystemName)
			}
		}

		c.Ctx.WriteString("ok")
	}
}

func (c *AdminCPController) Publish() {
	fun := c.Path(2)

	contentTypeId := c.PathInt64(3)
	contentType := database.DB.GetContentTypeWithId(contentTypeId)
	c.Data["ContentType"] = contentType

	if fun == "manage" {
		c.TplName = "admin/publish_manage.html"
		c.AddCSS("dataset.css")

		contents := database.DB.GetContentsWithContentTypeID(contentTypeId)
		c.Data["Contents"] = contents
	} else if fun == "create" {
		c.TplName = "admin/publish_create.html"
		c.AddCSS("dataset.css")
		c.AddCSS("datepicker.css")
		c.AddJS("ckeditor/ckeditor.js")
		c.AddJS("ckeditor/adapters/jquery.js")
		c.AddJS("date.js")
		c.AddJS("datePicker.js")

		languages := database.DB.GetLanguages()
		c.Data["Languages"] = languages

		fieldTypes := database.DB.GetFieldTypesWithContentTypeId(contentTypeId)
		c.Data["Fields"] = fieldTypes
		if c.IsPost() {
			title := c.GetString("title","")

			//权限
			privileges := c.GetStrings("privileges[]")
			privilegesJson, _ := utils.ToJSONWithSliceString(privileges)

			//发布语言
			languageId := c.GetString("language_id","")
			language := database.DB.GetLanguageWithLanguageID(utils.JKStrToInt64(languageId))

			//发布时间
			date := c.GetString("publish_date","")
			dateHour := c.GetString("publish_date_hour","")
			dateMinute := c.GetString("publish_date_minute","")
			dateAmpm := c.GetString("publish_date_ampm","")

			datetimeStr := utils.ValuesToDateTime(date, dateHour, dateMinute, dateAmpm)
			utils.JJKPrintln(datetimeStr)
			utils.JJKPrintln(date, dateHour, dateMinute, dateAmpm, languageId, privileges, title)

			//主表
			content := &models.Content{
				Language:language,
				ContentType:contentType,
				User:c.GetUser().(*models.User),
				ContentDate:time.Now(),
				ContentModified:time.Now(),
				ContentIsStandard:contentType.IsStandard,
				ContentTitle:title,
				ContentPrivileges:privilegesJson,
			}

			contentId, _ := database.DB.Orm.Insert(content)
			content.ContentId = contentId

			//自定义表
			params := utils.TemplateParams()
			params["FieldTypes"] = fieldTypes
			params["ContentType"] = contentType
			params["Content"] = content

			sql := utils.Template("INSERT INTO `{{.ContentType.SystemName}}` (`content_id`{{range .FieldTypes}},`{{.SystemName}}`{{end}}) VALUES ({{.Content.ContentId}}{{range .FieldTypes}},?{{end}});", params)
			var values []string
			for _, value := range fieldTypes {
				values = append(values, c.GetString(value.SystemName,""))
			}
			_, err := database.DB.DBBaseExecSQL(sql,values)
			if err != nil {
				utils.JJKPrintln(err)
			}

			c.SetError("添加成功", true)
		}
	} else if fun == "edit" {
		c.TplName = "admin/publish_edit.html"
		c.AddCSS("dataset.css")
		c.AddCSS("datepicker.css")
		c.AddJS("ckeditor/ckeditor.js")
		c.AddJS("ckeditor/adapters/jquery.js")
		c.AddJS("date.js")
		c.AddJS("datePicker.js")

		languages := database.DB.GetLanguages()
		c.Data["Languages"] = languages

		fieldTypes := database.DB.GetFieldTypesWithContentTypeId(contentTypeId)
		c.Data["Fields"] = fieldTypes

		contentId := c.PathInt64(4)
		content := database.DB.GetContentWithContentID(contentId)
		c.Data["Content"] = content

		params, _ := database.DB.DBBaseAnyTableSelectOneRowWithContentID(contentType.SystemName, contentId)
		utils.JJKPrintln(params)

		if c.IsPost() {
			title := c.GetString("title","")

			//权限
			privileges := c.GetStrings("privileges[]")
			privilegesJson, _ := utils.ToJSONWithSliceString(privileges)

			//发布语言
			languageId := c.GetString("language_id","")
			language := database.DB.GetLanguageWithLanguageID(utils.JKStrToInt64(languageId))

			//发布时间
			date := c.GetString("publish_date","")
			dateHour := c.GetString("publish_date_hour","")
			dateMinute := c.GetString("publish_date_minute","")
			dateAmpm := c.GetString("publish_date_ampm","")

			datetimeStr := utils.ValuesToDateTime(date, dateHour, dateMinute, dateAmpm)
			utils.JJKPrintln(datetimeStr)
			utils.JJKPrintln(date, dateHour, dateMinute, dateAmpm, languageId, privileges, title)

			content := &models.Content{
				Language:language,
				ContentType:contentType,
				User:c.GetUser().(*models.User),
				ContentDate:time.Now(),
				ContentModified:time.Now(),
				ContentIsStandard:contentType.IsStandard,
				ContentTitle:title,
				ContentPrivileges:privilegesJson,
			}

			contentId, _ := database.DB.Orm.Insert(content)
			content.ContentId = contentId

		}
	}
}