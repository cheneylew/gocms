package models

import "time"


//`orm:"-"`

type UserRole struct {
	UserRoleId int64	`orm:"pk;auto"`
	Name string
	Grade int64					//0 普通用户 1管理员 2超级管理员
}

type User struct {
	UserId int64		`orm:"pk;auto"`
	Username string
	Password string
	UserRole *UserRole	`orm:"rel(fk)"`
	UserFirstName string
	UserLastName string
	UserEmail string
	UserSalt string
	UserReferrer *User		`orm:"null;rel(fk)"`
	UserLastLogin time.Time		`orm:"null"`
	UserRegistDate time.Time	`orm:"null"`
	UserDeleted bool
	UserSex int 					//0 男 1女
}



