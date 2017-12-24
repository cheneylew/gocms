package models

import "time"

type FieldType struct {
	FieldTypeId int64 `orm:"pk;auto"`
	ContentType *ContentType `orm:"rel(fk)"`
	SystemName string
	Name string
	Help string
	Required bool
	DefaultValue string
	Options string			//一些选项
}

type ContentType struct {
	ContentTypeId int64 `orm:"pk;auto"`
	SystemName string
	Name string
	IsStandard bool
	IsPrivileged bool
	Template string
	BaseUrl string
	FieldTypes []*FieldType `orm:"reverse(many)"`
}

type Content struct {
	ContentId int64 `orm:"pk;auto"`
	Language *Language `orm:"rel(fk)"`
	ContentType *ContentType `orm:"rel(fk)"`
	User *User `orm:"rel(fk)"`
	ContentDate time.Time
	ContentModified time.Time
	ContentTopics string
	ContentIsStandard bool
	ContentTitle string
	ContentPrivileges string
	ContentHits string
}

type Language struct {
	LanguageId int64 `orm:"pk;auto"`
	Name string
}