package models

import (
	"time"
	"fmt"
	"github.com/cheneylew/goutil/utils"
)

type FieldType struct {
	FieldTypeId int64 `orm:"pk;auto"`
	ContentType *ContentType `orm:"rel(fk)"`
	SystemName string
	Name string
	Type string
	Help string
	Required bool
	DefaultValue string
	Options string			//一些选项
}

func (f *FieldType)RequiredHTML() string {
	return fmt.Sprintf(`<li>
	<label for="required">Required Field</label>
	<input type="checkbox" name="required" value="%d" class="checkbox" />
	<div class="help">If checked, this box must be checked for the form to be processed.</div>
	</li>`, f.Required)
}

func (f *FieldType)HelpHTML() string {
	return fmt.Sprintf(`<li>
	<label for="help">Help Text</label>
	<textarea name="help" style="width: 500px; height: 80px" class="textarea">%s</textarea>
	<div class="help">This help text will be displayed beneath the field.  Use it to guide the user in responding correctly.</div>
	</li>`, f.Help)
}

func (f *FieldType)ToInputHTML() string {
	params := utils.TemplateParams()
	params["FieldType"] = f

	tplStr := ""
	if f.Type == "checkbox" {
		tplStr = `<li id="row_{{.FieldType.FieldTypeId}}">
			<label for="{{.FieldType.SystemName}}">{{.FieldType.SystemName}}</label>
			<input type="{{.FieldType.Type}}" name="{{.FieldType.SystemName}}" value="{{.FieldType.DefaultValue}}" class="checkbox" {{if eq .FieldType.DefaultValue "1"}}checked="checked"{{end}}>
			</li>`
	}

	return utils.Template(tplStr, params)
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