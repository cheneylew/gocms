package models

import "github.com/astaxie/beego/orm"

func init() {
	orm.RegisterModel(
		new(User),
		new(UserRole),
		new(FieldType),
		new(ContentType),
		new(Language),
		new(Content),
		new(UserLogins),
	)
}