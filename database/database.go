package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/cheneylew/gocms/models"
	"github.com/cheneylew/goutil/utils/beego"
	"github.com/astaxie/beego/orm"
)

var DB DataBase

func init() {
	db := beego.InitRegistDB("cheneylew","12344321","47.91.151.207","3308","gocms")
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