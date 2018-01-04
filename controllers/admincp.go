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
	"github.com/cheneylew/goutil/utils/beego"
	"path"
	"os"
	"encoding/json"
	"github.com/astaxie/beego/orm"
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
				Name:name,
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
				c.RedirectWithURL("/admincp/types?active=2")
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

		rules := make(map[string]interface{}, 0)
		if fieldType == models.FieldTypeStrCheckbox {
			database.DB.DBBaseAddColumnTinyInt(tableInfo.SystemName, systemName)
		} else if fieldType == models.FieldTypeStrWysiwyg {
			err := database.DB.DBBaseAddColumnVarChar(tableInfo.SystemName, systemName, 2000)
			if err != nil {
				panic(err)
			}
		} else if fieldType == models.FieldTypeStrDate {

			rules["future_only"] = c.GetString("future_only")
			err := database.DB.DBBaseAddColumn(tableInfo.SystemName, systemName, beego.DataTypeDate, "0000-00-00")
			if err != nil {
				panic(err)
			}
		} else if fieldType == models.FieldTypeStrDatetime {

			rules["future_only"] = c.GetString("future_only")
			err := database.DB.DBBaseAddColumn(tableInfo.SystemName, systemName, beego.DataTypeDateTime, "0000-00-00 00:00:00")
			if err != nil {
				panic(err)
			}
		} else if fieldType == models.FieldTypeStrTextarea {

			rules["width"] = c.GetString("width","")
			rules["height"] = c.GetString("height","")
			rules["validators"] = c.GetStrings("validators[]")

			err := database.DB.DBBaseAddColumnVarChar(tableInfo.SystemName, systemName, 5000)
			if err != nil {
				panic(err)
			}
		} else if fieldType == models.FieldTypeStrText {

			rules["width"] = c.GetString("width","")
			rules["validators"] = c.GetStrings("validators[]")

			err := database.DB.DBBaseAddColumnVarChar(tableInfo.SystemName, systemName, 5000)
			if err != nil {
				panic(err)
			}
		}  else if fieldType == models.FieldTypeStrMulticheckbox|| fieldType == models.FieldTypeStrMultiselect || fieldType == models.FieldTypeStrSelect || fieldType == models.FieldTypeStrRadio{

			model.Options = utils.JKJSON(helper.StringToOptions(c.GetString("options","")))
			model.DefaultValue = strings.Join(helper.StringToOptionValue(model.DefaultValue), "|")

			err := database.DB.DBBaseAddColumnVarChar255(tableInfo.SystemName, systemName)
			if err != nil {
				panic(err)
			}
		} else  if fieldType == models.FieldTypeStrMemberGroupRelationship {
			rules["allowMultiple"] = c.GetString("allow_multiple","")

			err := database.DB.DBBaseAddColumnVarChar255(tableInfo.SystemName, systemName)
			if err != nil {
				panic(err)
			}
		} else  if fieldType == models.FieldTypeStrFileUpload {
			rules["width"] = c.GetString("width","")
			rules["filetypes"] = c.GetStrings("filetypes")

			err := database.DB.DBBaseAddColumnVarChar255(tableInfo.SystemName, systemName)
			if err != nil {
				panic(err)
			}
		} else if fieldType == models.FieldTypeStrRelationship {
			rules["fieldName"] = c.GetString("field_name","")
			rules["contentType"] = c.GetString("content_type","")
			rules["allowMultiple"] = c.GetString("allow_multiple","")

			err := database.DB.DBBaseAddColumnVarChar255(tableInfo.SystemName, systemName)
			if err != nil {
				panic(err)
			}
		} else {
			panic("have not setting...")
			return
		}

		model.RuleJson = utils.JKJSON(rules)
		_, err := database.DB.Orm.Insert(model)
		if err != nil {
			utils.JJKPrintln(err)
		} else {
			utils.JJKPrintln("field add ok!")
		}
	}
}

func (c *AdminCPController) CustomFields() {
	fun := c.Path(2)
	if fun == "ajax_field_form" {
		ftype := c.GetString("type","")
		fieldType := models.FieldType{
			Type:ftype,
		}
		html := fieldType.RuleHTML()
		params := utils.TemplateParams()
		if ftype == models.FieldTypeStrMemberGroupRelationship {
			params["Roles"] = database.DB.GetUserRoles()
		} else if ftype == models.FieldTypeStrRelationship {
			params["ContentTypes"] = database.DB.GetContentTypes()
		}
		c.Ctx.WriteString(utils.Template(html,params))
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
		info := c.GetString("info")
		for _, value := range itemIds {
			if action == "delete_articles_not_standard" {
				contentType := database.DB.GetContentTypeWithId(utils.JKStrToInt64(info))
				//删除自定义表行
				database.DB.DBBaseDeleteRowWithId(contentType.SystemName, value)
			} else if action == "delete_articles" {
				//删除Content
				item := database.DB.GetContentWithContentID(value)
				database.DB.Orm.Delete(item)

				//删除自定义表行
				database.DB.DBBaseDeleteRowWithContentId(item.ContentType.SystemName, value)
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

		if contentType.IsStandard {
			contents := database.DB.GetContentsWithContentTypeIDAndOffsetLimit(contentTypeId, c.GetOffset(), c.GetLimit())
			c.Data["Contents"] = contents
		} else {
			var list []orm.Params
			sql := fmt.Sprintf("select * from %s order by %s_id desc limit %d, %d;", contentType.SystemName, contentType.SystemName, c.GetOffset(), c.GetLimit())
			database.DB.Orm.Raw(sql).Values(&list)
			c.Data["ListMaps"] = list
			c.Data["RowIDStr"] = fmt.Sprintf("%s_id", contentType.SystemName)

			fieldTypes := database.DB.GetFieldTypesWithContentTypeId(contentTypeId)
			maxCount := 5
			if len(fieldTypes) < 5 {
				maxCount = len(fieldTypes)
			}
			c.Data["FieldTypes"] = fieldTypes[:maxCount]
		}

		count := database.DB.DBBaseTableCount(contentType.SystemName)
		c.Data["RowCount"] = count

		//分页
		c.Data["Pagination"] = c.Pagination(count, conf.GlobalConfig.PageLimit)
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
		for i := 0; i < len(fieldTypes); i++ {
			if fieldTypes[i].Type == models.FieldTypeStrMemberGroupRelationship {
				fieldTypes[i].Options = utils.JKJSON(database.DB.GetUserRoles())
			} else if fieldTypes[i].Type == models.FieldTypeStrRelationship {
				var rule models.FieldTypeTextViewRule
				json.Unmarshal([]byte(fieldTypes[i].RuleJson), &rule)
				contentType := database.DB.GetContentTypeWithId(utils.JKStrToInt64(rule.ContentType))
				contentTypeName := contentType.SystemName
				fieldName := rule.FieldName

				var list []orm.Params
				sql := fmt.Sprintf("select %s_id, %s from %s", contentTypeName, fieldName, contentTypeName)
				utils.JJKPrintln(sql)
				database.DB.Orm.Raw(sql).Values(&list)

				fieldTypes[i].Options = utils.JKJSON(list)
				fieldTypes[i].Other = contentTypeName+"_id"+","+fieldName
			}
		}

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

			//发布内容表
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

			if contentType.IsStandard {
				contentId, _ := database.DB.Orm.Insert(content)
				content.ContentId = contentId
			} else {
				content.ContentId = 0
			}


			//自定义字段表
			params := utils.TemplateParams()
			params["FieldTypes"] = fieldTypes
			params["ContentType"] = contentType
			params["Content"] = content

			sql := utils.Template("INSERT INTO `{{.ContentType.SystemName}}` (`content_id`{{range .FieldTypes}},`{{.SystemName}}`{{end}}) VALUES ({{.Content.ContentId}}{{range .FieldTypes}},?{{end}});", params)
			var values []interface{}
			for _, value := range fieldTypes {
				if value.Type == models.FieldTypeStrCheckbox {
					if c.GetString(value.SystemName,"") == "on" {
						values = append(values, 1)
					} else {
						values = append(values, 0)
					}
				} else if value.Type == models.FieldTypeStrDatetime {
					datetimeStr := c.GetDateTimeStr(value.SystemName)
					values = append(values, datetimeStr)

				} else if value.Type == models.FieldTypeStrMulticheckbox || value.Type == models.FieldTypeStrMultiselect {
					values = append(values, strings.Join(c.GetStrings(value.SystemName+"[]"), "|"))
				} else if value.Type == models.FieldTypeStrFileUpload {
					lastFile := c.GetString(value.SystemName+"_uploaded")
					deleteFileFileUpload, _ := c.GetBool("delete_file_file_upload", false)

					f, h, err := c.GetFile(value.SystemName)
					if f != nil && err == nil {
						defer f.Close()
						newFileName := utils.MD5(utils.RandomString(32)+time.Now().String())+utils.Ext(h.Filename)
						newPath := "static/upload/" + newFileName
						c.SaveToFile(value.SystemName, newPath)

						//生成预览文件
						realPath := utils.SelfDir()+"/"+newPath
						utils.ImageThumbnail(realPath, 150)

						values = append(values, newPath)
					} else {
						if deleteFileFileUpload {
							values = append(values, "")
						} else {
							values = append(values, strings.TrimLeft(lastFile,"/"))
						}
					}

					//删除上次上传的照片
					lastFileThumbnail := utils.ThumbnailPath(lastFile)
					if (len(lastFile) > 0 && f != nil) || deleteFileFileUpload {
						//本次传了文件，移除上次文件
						lastFile = path.Join(utils.SelfDir(), lastFile)
						lastFileThumbnail = path.Join(utils.SelfDir(), lastFileThumbnail)
						os.Remove(lastFile)
						os.Remove(lastFileThumbnail)
					}
				} else if value.Type == models.FieldTypeStrMemberGroupRelationship || value.Type == models.FieldTypeStrRelationship {
					var rule models.FieldTypeTextViewRule
					json.Unmarshal([]byte(value.RuleJson), &rule)
					if rule.AllowMultiple == "1" {
						values = append(values, strings.Join(c.GetStrings(value.SystemName+"[]"),"|"))
					} else {
						values = append(values, c.GetString(value.SystemName,""))
					}
				} else {
					values = append(values, c.GetString(value.SystemName,""))
				}
			}

			_, err := database.DB.DBBaseExecSQL(sql,values)
			if err != nil {
				utils.JJKPrintln(err)
			}

			c.SetError("添加成功", true)
			c.RedirectWithURL(fmt.Sprintf("/admincp/publish/manage/%d?active=2", contentType.ContentTypeId))
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

		c.Data["ContentType"] = contentType

		var rowId int64
		var content *models.Content
		if contentType.IsStandard {
			contentId := c.PathInt64(4)
			content = database.DB.GetContentWithContentID(contentId)
			c.Data["Content"] = content
			rowId = contentId
		} else {
			rowId = c.PathInt64(4)
		}

		fieldTypes := database.DB.GetFieldTypesWithContentTypeId(contentTypeId)

		var params orm.Params
		if !contentType.IsStandard {
			params, _ = database.DB.DBBaseAnyTableSelectOneRowWithID(contentType.SystemName, rowId)
		} else {
			params, _ = database.DB.DBBaseAnyTableSelectOneRowWithContentID(contentType.SystemName, rowId)
		}
		for i := 0; i < len(fieldTypes); i++ {
			fieldTypes[i].DefaultValue = utils.ToString(params[fieldTypes[i].SystemName])

			//回填处理
			if fieldTypes[i].Type == models.FieldTypeStrMemberGroupRelationship {
				fieldTypes[i].Options = utils.JKJSON(database.DB.GetUserRoles())
			} else if fieldTypes[i].Type == models.FieldTypeStrRelationship {
				var rule models.FieldTypeTextViewRule
				json.Unmarshal([]byte(fieldTypes[i].RuleJson), &rule)
				contentType := database.DB.GetContentTypeWithId(utils.JKStrToInt64(rule.ContentType))
				contentTypeName := contentType.SystemName
				fieldName := rule.FieldName

				var list []orm.Params
				sql := fmt.Sprintf("select %s_id, %s from %s", contentTypeName, fieldName, contentTypeName)
				utils.JJKPrintln(sql)
				database.DB.Orm.Raw(sql).Values(&list)

				fieldTypes[i].Options = utils.JKJSON(list)
				fieldTypes[i].Other = contentTypeName+"_id"+","+fieldName
			}
		}

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
			datetimeStr := c.GetDateTime("publish_date")
			utils.JJKPrintln(datetimeStr)
			utils.JJKPrintln(languageId, privileges, title)

			content.Language = language
			content.ContentType = contentType
			content.User = c.GetUser().(*models.User)
			content.ContentDate = time.Now()
			content.ContentModified = time.Now()
			content.ContentIsStandard = contentType.IsStandard
			content.ContentTitle = title
			content.ContentPrivileges = privilegesJson;

			if contentType.IsStandard {
				database.DB.Orm.Update(content)
			} else {
				content.ContentId = 0
			}

			//自定义表
			params := utils.TemplateParams()
			params["FieldTypes"] = fieldTypes
			params["ContentType"] = contentType
			params["Content"] = content

			sql := utils.Template("UPDATE `{{.ContentType.SystemName}}` SET `content_id` = {{.Content.ContentId}}{{range .FieldTypes}},`{{.SystemName}}` =? {{end}} WHERE `content_id` = {{.Content.ContentId}};", params)
			var values []interface{}
			for _, value := range fieldTypes {
				if value.Type == models.FieldTypeStrCheckbox {
					if c.GetString(value.SystemName,"") == "on" {
						values = append(values, 1)
					} else {
						values = append(values, 0)
					}
				} else if value.Type == models.FieldTypeStrDatetime {
					datetimeStr := c.GetDateTimeStr(value.SystemName)
					values = append(values, datetimeStr)
				} else if value.Type == models.FieldTypeStrMulticheckbox || value.Type == models.FieldTypeStrMultiselect {
					values = append(values, strings.Join(c.GetStrings(value.SystemName+"[]"), "|"))
				} else if value.Type == models.FieldTypeStrFileUpload {
					lastFile := c.GetString(value.SystemName+"_uploaded")
					deleteFileFileUpload, _ := c.GetBool("delete_file_file_upload", false)

					f, h, err := c.GetFile(value.SystemName)
					if f != nil && err == nil {
						defer f.Close()
						newFileName := utils.MD5(utils.RandomString(32)+time.Now().String())+utils.Ext(h.Filename)
						newPath := "static/upload/" + newFileName
						c.SaveToFile(value.SystemName, newPath)

						//生成预览文件
						realPath := utils.SelfDir()+"/"+newPath
						utils.ImageThumbnail(realPath, 150)

						values = append(values, newPath)
					} else {
						if deleteFileFileUpload {
							values = append(values, "")
						} else {
							values = append(values, strings.TrimLeft(lastFile,"/"))
						}
					}

					//删除上次上传的照片
					lastFileThumbnail := utils.ThumbnailPath(lastFile)
					if (len(lastFile) > 0 && f != nil) || deleteFileFileUpload {
						//本次传了文件，移除上次文件
						lastFile = path.Join(utils.SelfDir(), lastFile)
						lastFileThumbnail = path.Join(utils.SelfDir(), lastFileThumbnail)
						os.Remove(lastFile)
						os.Remove(lastFileThumbnail)
					}
				} else if value.Type == models.FieldTypeStrMemberGroupRelationship || value.Type == models.FieldTypeStrRelationship{
					var rule models.FieldTypeTextViewRule
					json.Unmarshal([]byte(value.RuleJson), &rule)
					if rule.AllowMultiple == "1" {
						values = append(values, strings.Join(c.GetStrings(value.SystemName+"[]"),"|"))
					} else {
						values = append(values, c.GetString(value.SystemName,""))
					}
				} else {
					values = append(values, c.GetString(value.SystemName,""))
				}

			}
			utils.JJKPrintln(values)
			utils.JJKPrintln(sql)
			_, err := database.DB.DBBaseExecSQL(sql,values)
			if err != nil {
				utils.JJKPrintln(err)
			}

			c.RedirectWithURL(c.Ctx.Request.RequestURI)
		}
	}
}