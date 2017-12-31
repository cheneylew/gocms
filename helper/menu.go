package helper

import (
	"github.com/cheneylew/gocms/database"
	"fmt"
	"github.com/cheneylew/gocms/models"
)

func CreateMenus() []*models.Menu {
	var menus []*models.Menu
	

	item := &models.Menu{
		MenuID:1,
		ParentId:0,
		Name:"仪表盘",
		Class:"dashboard",
		Icon:"",
		Url:"/admincp/home",
		IsActive:true,
		Chirldren:nil,
	}
	menus = append(menus, item)

	types := database.DB.GetContentTypes()
	var publishMenus []*models.Menu
	publishMenus = append(publishMenus,
		&models.Menu{
			MenuID:10,
			ParentId:2,
			Name:"发布文章",
			Icon:"",
			Url:"/admincp/home",
			IsActive:true,
			Chirldren:nil,
		})
	for _, value := range types {
		publishMenus = append(publishMenus,
			&models.Menu{
			MenuID:11,
			ParentId:2,
			Name:fmt.Sprintf("管理%s", value.Name),
			Icon:"",
			Url:fmt.Sprintf("/admincp/publish/manage/%d", value.ContentTypeId),
			IsActive:true,
			Chirldren:nil,
		})
	}
	publishMenus = append(publishMenus,
		&models.Menu{
			MenuID:11,
			ParentId:2,
			Name:"内容模型",
			Icon:"",
			Url:"/admincp/types",
			IsActive:true,
			Chirldren:nil,
		})

	item = &models.Menu{
		MenuID:2,
		ParentId:0,
		Name:"文章管理",
		Class:"publish",
		Icon:"",
		Url:"/admincp/home",
		IsActive:false,
		Chirldren:publishMenus,
	}
	menus = append(menus, item)

	item = &models.Menu{
		MenuID:3,
		ParentId:0,
		Name:"会员",
		Class:"members",
		Icon:"",
		Url:"/admincp/home",
		IsActive:false,
		Chirldren:[]*models.Menu{
			&models.Menu{
				MenuID:31,
				ParentId:3,
				Name:"管理用户",
				Icon:"",
				Url:"/admincp/users",
				IsActive:true,
				Chirldren:nil,
			},
			&models.Menu{
				MenuID:31,
				ParentId:3,
				Name:"添加用户",
				Icon:"",
				Url:"/admincp/users/add",
				IsActive:true,
				Chirldren:nil,
			},
		},
	}
	menus = append(menus, item)

	item = &models.Menu{
		MenuID:4,
		ParentId:0,
		Name:"报告",
		Class:"reports",
		Icon:"",
		Url:"/admincp/home",
		IsActive:false,
		Chirldren:nil,
	}
	menus = append(menus, item)

	item = &models.Menu{
		MenuID:5,
		ParentId:0,
		Name:"设计",
		Class:"design",
		Icon:"",
		Url:"/admincp/home",
		IsActive:false,
		Chirldren:nil,
	}
	menus = append(menus, item)

	item = &models.Menu{
		MenuID:6,
		ParentId:0,
		Name:"配置",
		Class:"configuration",
		Icon:"",
		Url:"/admincp/home",
		IsActive:false,
		Chirldren:nil,
	}
	menus = append(menus, item)

	return menus
}

