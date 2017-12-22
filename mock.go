package main

import (
	"github.com/cheneylew/gocms/database"
	"github.com/cheneylew/gocms/models"
	"github.com/cheneylew/goutil/utils"
	"time"
)

func MockMain()  {
	roleCount := database.DB.DBBaseTableCount("user_role")
	if roleCount == 0 {
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

	userCount := database.DB.DBBaseTableCount("user")
	if userCount == 0 {
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
}
