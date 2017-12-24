package main

import (
	"github.com/cheneylew/gocms/database"
	"github.com/cheneylew/gocms/models"
	"github.com/cheneylew/goutil/utils"
	"time"
)

func MockMain()  {
	if database.DB.DBBaseTableCount("user_role") == 0 {
		var role *models.UserRole
		role = &models.UserRole{
			Name:"普通用户",
			Grade:0,
		}
		database.DB.Orm.Insert(role)

		role = &models.UserRole{
			Name:"管理员",
			Grade:1,
		}
		database.DB.Orm.Insert(role)

		role = &models.UserRole{
			Name:"超级管理员",
			Grade:2,
		}
		database.DB.Orm.Insert(role)
	}

	if database.DB.DBBaseTableCount("user") == 0 {
		var obj *models.User
		obj = &models.User{
			Username:"cheneylew",
			Password:"tough1988",
			UserLastLogin:time.Now(),
			UserReferrer:&models.User{
				UserId:0,
			},
			UserRole:&models.UserRole{
				UserRoleId:3,
				Name:"超级管理员",
				Grade:2,
			},
		}
		a, err := database.DB.Orm.Insert(obj)
		utils.JJKPrintln(a, err)
	}

	if database.DB.DBBaseTableCount("language") == 0 {
		database.DB.Orm.InsertMulti(3, []*models.Language{&models.Language{Name:"简体中文"},&models.Language{Name:"繁体中文"},&models.Language{Name:"English"}})
	}

	//ctype := &models.ContentType{
	//	SystemName:"book",
	//	Name:"Book",
	//	IsStandard:false,
	//}
	//database.DB.Orm.Insert(ctype)
}
