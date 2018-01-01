package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/cheneylew/gocms/models"
	"github.com/cheneylew/goutil/utils/beego"
	"github.com/astaxie/beego/orm"
)

var DB DataBase

func init() {
	//db := beego.InitRegistDB("cheneylew","12344321","47.91.151.207","3308","gocms")
	db := beego.InitRegistDB("root","cnldj1988","127.0.0.1","3306","gocms")
	DB = DataBase{
		BaseDataBase:*db,
	}
}

type DataBase struct {
	beego.BaseDataBase
}

func (db *DataBase)GetUsersWithEmailOrUsername(emailOrUsername string) []*models.User {
	var models []*models.User
	qs := db.Orm.QueryTable("user")

	cond := orm.NewCondition().And("Username", emailOrUsername).Or("UserEmail", emailOrUsername)
	_, err := qs.SetCond(cond).RelatedSel("UserRole").All(&models)
	if err != nil {
		return nil
	}

	return models
}

func (db *DataBase)GetUserWithEmailOrUsername(emailOrUsername string) *models.User {
	var models []*models.User
	qs := db.Orm.QueryTable("user")

	cond := orm.NewCondition().And("Username", emailOrUsername).Or("UserEmail", emailOrUsername)
	_, err := qs.SetCond(cond).RelatedSel("UserRole").All(&models)
	if err != nil || len(models) == 0 {
		return nil
	}

	return models[0]
}

func (db *DataBase)GetUsers() []*models.User {
	var models []*models.User
	qs := db.Orm.QueryTable("User")

	_, err := qs.RelatedSel("UserRole").All(&models)
	if err != nil {
		return nil
	}

	return models
}

func (db *DataBase)GetUserLogins() []*models.UserLogins {
	var models []*models.UserLogins
	qs := db.Orm.QueryTable("UserLogins")

	_, err := qs.RelatedSel("User").Limit(30,0).All(&models)
	if err != nil {
		return nil
	}

	return models
}

func (db *DataBase)GetContentTypes() []*models.ContentType {
	var models []*models.ContentType
	qs := db.Orm.QueryTable("ContentType")

	_, err := qs.All(&models)
	if err != nil {
		return nil
	}

	return models
}

func (db *DataBase)GetContentTypeWithId(contentTypeId int64) *models.ContentType {
	var models []*models.ContentType
	qs := db.Orm.QueryTable("ContentType")

	a, err := qs.Filter("ContentTypeId", contentTypeId).All(&models)
	if err != nil || a == 0 {
		return nil
	}

	return models[0]
}

func (db *DataBase)GetUserRoles() []*models.UserRole {
	var models []*models.UserRole
	qs := db.Orm.QueryTable("UserRole")

	_, err := qs.All(&models)
	if err != nil {
		return nil
	}

	return models
}

func (db *DataBase)GetUserRolesWithGrade(grade int64) *models.UserRole {
	var models []*models.UserRole
	qs := db.Orm.QueryTable("UserRole")

	_, err := qs.Filter("Grade", grade).All(&models)
	if err != nil || len(models) == 0{
		return nil
	}

	return models[0]
}

func (db *DataBase)GetFieldTypesWithContentTypeId(contentTypeId int64) []*models.FieldType {
	var models []*models.FieldType
	qs := db.Orm.QueryTable("FieldType")

	_, err := qs.Filter("ContentType__ContentTypeId", contentTypeId).All(&models)
	if err != nil || len(models) == 0{
		return nil
	}

	return models
}

func (db *DataBase)GetFieldTypesWithFieldTypeID(fieldTypeId int64) *models.FieldType {
	var models []*models.FieldType
	qs := db.Orm.QueryTable("FieldType")

	_, err := qs.Filter("FieldTypeId", fieldTypeId).RelatedSel().All(&models)
	if err != nil || len(models) == 0{
		return nil
	}

	return models[0]
}

func (db *DataBase)GetLanguages() []*models.Language {
	var models []*models.Language
	qs := db.Orm.QueryTable("Language")

	_, err := qs.All(&models)
	if err != nil {
		return nil
	}

	return models
}

func (db *DataBase)GetLanguageWithLanguageID(languageId int64) *models.Language {
	var models []*models.Language
	qs := db.Orm.QueryTable("Language")

	_, err := qs.Filter("LanguageId", languageId).RelatedSel().All(&models)
	if err != nil || len(models) == 0{
		return nil
	}

	return models[0]
}


func (db *DataBase)GetContents() []*models.Content {

	var models []*models.Content
	qs := db.Orm.QueryTable("Content")

	_, err := qs.RelatedSel().All(&models)
	if err != nil {
		return nil
	}

	return models
}

func (db *DataBase)GetContentWithContentID(contentId int64) *models.Content {
	var models []*models.Content
	qs := db.Orm.QueryTable("Content")

	_, err := qs.Filter("ContentId", contentId).RelatedSel().All(&models)
	if err != nil || len(models) == 0{
		return nil
	}

	return models[0]
}

func (db *DataBase)GetContentsWithContentTypeID(contentTypeId int64) []*models.Content {
	var models []*models.Content
	qs := db.Orm.QueryTable("Content")

	_, err := qs.Filter("ContentType__ContentTypeId", contentTypeId).RelatedSel().All(&models)
	if err != nil || len(models) == 0{
		return nil
	}

	return models
}