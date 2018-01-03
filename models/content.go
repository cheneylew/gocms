package models

import (
	"time"
	"fmt"
	"github.com/cheneylew/goutil/utils"
	"encoding/json"
	"strings"
)

func StringToOptionValue(defaultValue string) []string {
	tmpValues := strings.Split(defaultValue, "\n")
	var values []string
	for _, value := range tmpValues {
		value = utils.Trim(value)
		if len(value) > 0 {
			values = append(values, value)
		}
	}

	return values
}

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
	RuleJson string			//规则限制
	Other string
}

type Option struct {
	Key string
	Value string
	IsSelected bool
}

type FieldTypeTextViewRule struct {
	Height     string   `json:"height"`
	Validators []string `json:"validators"`
	Width      string   `json:"width"`
	ContentType string	`json:"contentType"`
	FieldName string	`json:"fieldName"`
	AllowMultiple string `json:"allowMultiple`
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
	header := `<li id="row_{{.FieldType.FieldTypeId}}">
						<label for="{{.FieldType.SystemName}}">{{.FieldType.SystemName}}</label>`
	footer := `<div class="help">{{.FieldType.Help}}</div></li>`
	if f.Type == FieldTypeStrCheckbox {
		tplStr = `<input type="{{.FieldType.Type}}" name="{{.FieldType.SystemName}}" {{if Equal .FieldType.DefaultValue "1"}}checked="checked"{{end}} class="checkbox" {{if eq .FieldType.DefaultValue "1"}}checked="checked"{{end}}>`
	} else if f.Type == FieldTypeStrWysiwyg {
		tplStr = `<div style="float: left; width:  750px ">
										<textarea type="textarea" id="{{.FieldType.SystemName}}" name="{{.FieldType.SystemName}}" style="width: 700px; height: 140px " id="body" class="basic wysiwyg">
										{{.FieldType.DefaultValue}}
										</textarea>
									</div>`
	} else if f.Type == FieldTypeStrDate {
		tplStr = `<input type="text" name="{{.FieldType.SystemName}}" id="{{.FieldType.SystemName}}" value="{{.FieldType.DefaultValue}}" style="width: 80px" class="text date datepick dp-applied"><a href="#" class="dp-choose-date" title="Choose date">Choose date</a>`
	} else if f.Type == FieldTypeStrDatetime {
		var hours []string
		var minutes []string
		for i := 0; i<= 12; i++ {
			hours = append(hours, fmt.Sprintf("%02d", i))
		}
		for i := 0; i<= 59; i++ {
			minutes = append(minutes, fmt.Sprintf("%02d", i))
		}
		params["Hours"] = hours
		params["Minutes"] = minutes

		params["DefaultDate"] = ""
		params["DefaultHour"] = ""
		params["DefaultAMPM"] = ""
		params["DefaultMinute"] = ""

		if len(f.DefaultValue) > 2 {
			defaultDateTime := utils.JKStringToTime(f.DefaultValue)
			params["DefaultDate"] = utils.JKDateToString(defaultDateTime)
			if defaultDateTime.Hour() > 12 {
				params["DefaultHour"] = fmt.Sprintf("%02d", defaultDateTime.Hour() - 12)
				params["DefaultAMPM"] = "pm"
			} else {
				params["DefaultHour"] = fmt.Sprintf("%02d", defaultDateTime.Hour())
				params["DefaultAMPM"] = "am"
			}
			params["DefaultMinute"] = fmt.Sprintf("%02d", defaultDateTime.Minute())
		}

		tplStr = `<input type="text" name="{{.FieldType.SystemName}}" id="{{.FieldType.SystemName}}" value="{{.DefaultDate}}" style="width: 80px" class="text datetime datepick dp-applied"><a href="#" class="dp-choose-date" title="Choose date">Choose date</a>
		<select name="{{.FieldType.SystemName}}_hour">
		{{range .Hours}}
		{{if eq . $.DefaultHour}}
		<option value="{{.}}" selected="selected">{{.}}</option>
		{{else}}
		<option value="{{.}}">{{.}}</option>
		{{end}}
		{{end}}
		</select>
		<select name="{{.FieldType.SystemName}}_minute">
		{{range .Minutes}}
		{{if eq . $.DefaultMinute}}
		<option value="{{.}}" selected="selected">{{.}}</option>
		{{else}}
		<option value="{{.}}">{{.}}</option>
		{{end}}
		{{end}}
		</select>
		<select name="{{.FieldType.SystemName}}_ampm">
			<option value="am" {{if Equal .DefaultAMPM "am"}}selected="selected"{{end}}>am</option>
			<option value="pm" {{if Equal .DefaultAMPM "pm"}}selected="selected"{{end}}>pm</option>
		</select>`
	} else if f.Type == FieldTypeStrTextarea {
		var rule FieldTypeTextViewRule
		err := json.Unmarshal([]byte(f.RuleJson), &rule)
		if err == nil {
			params["Rule"] = rule
		}
		tplStr = `<textarea name="{{.FieldType.SystemName}}" style="width: {{.Rule.Width}}; height: {{.Rule.Height}}" class="{{if .FieldType.Required}}required {{end}} textarea {{if .Rule}}{{range .Rule.Validators}}{{.}} {{end}}{{end}}">{{.FieldType.DefaultValue}}</textarea>
				 `
	} else if f.Type == FieldTypeStrText {
		var rule FieldTypeTextViewRule
		err := json.Unmarshal([]byte(f.RuleJson), &rule)
		if err == nil {
			params["Rule"] = rule
		}
		tplStr = `<input type="text" name="{{.FieldType.SystemName}}" id="{{.FieldType.SystemName}}" value="{{.FieldType.DefaultValue}}" style="width: {{.Rule.Width}}" class="{{if .FieldType.Required}}required {{end}} text {{if .Rule}}{{range .Rule.Validators}}{{.}} {{end}}{{end}}">`
	} else if f.Type == FieldTypeStrMulticheckbox || f.Type == FieldTypeStrMultiselect || f.Type == FieldTypeStrSelect || f.Type == FieldTypeStrRadio{
		var options []Option
		err := json.Unmarshal([]byte(f.Options), &options)
		values := strings.Split(f.DefaultValue, "|")
		if err == nil {
			for i:=0; i< len(options); i++ {
				options[i].IsSelected = utils.InSlice(options[i].Value, values)
			}
			params["Options"] = options
			utils.JJKPrintln(options, values, f.DefaultValue)
		}
		if f.Type == FieldTypeStrMulticheckbox {
			tplStr = `<div style="float: left">
			{{range .Options}}
			<div class="check_option">
				<input type="checkbox" name="{{$.FieldType.SystemName}}[]" {{if .IsSelected}}checked="checked"{{end}} value="{{.Value}}">{{.Key}}
			</div>
			{{end}}
			</div>`
		} else if f.Type == FieldTypeStrMultiselect || f.Type == FieldTypeStrSelect {
			params["IsMulti"] = f.Type == FieldTypeStrMultiselect
			tplStr = `<div style="float: left">
			<select name="{{$.FieldType.SystemName}}{{if .IsMulti}}[]{{end}}" {{if .IsMulti}}multiple="multiple"{{end}}>
			{{range .Options}}
				<option {{if .IsSelected}}selected="selected"{{end}} value="{{.Value}}">{{.Key}}</option>
			{{end}}
			</select>
			</div>`
		} else if f.Type == FieldTypeStrRadio {
			tplStr = `
			{{range .Options}}
			<input type="radio" name="{{$.FieldType.SystemName}}" {{if .IsSelected}}checked="checked"{{end}} class="{{if .FieldType.Required}}required {{end}} radio" value="{{.Value}}"/> {{.Key}}&nbsp;&nbsp;&nbsp;
			{{end}}
			`
		}
	} else if f.Type == FieldTypeStrFileUpload {
		var rule FieldTypeTextViewRule
		err := json.Unmarshal([]byte(f.RuleJson), &rule)
		if err == nil {
			params["Rule"] = rule
			params["ThumbnailPath"] = utils.ThumbnailPath(f.DefaultValue)
		}
		tplStr = `<input type="hidden" name="{{$.FieldType.SystemName}}_uploaded" value="/{{.FieldType.DefaultValue}}">
						<input type="hidden" name="{{$.FieldType.SystemName}}_ftp" value="">
						<input id="{{$.FieldType.SystemName}}" type="file" name="{{$.FieldType.SystemName}}" style="width: {{.Rule.Width}}" class="{{if .FieldType.Required}}required {{end}}file text">
						{{if .FieldType.DefaultValue}}
						<br>
<a href="/{{.FieldType.DefaultValue}}" target="_blank">
<img style="margin-left: 150px" src="/{{.ThumbnailPath}}" alt="preview">
</a><br>
<input style="margin-left: 150px" type="checkbox" name="delete_file_file_upload" value="1"> <span style="font-style: italic; color: #ff6464">Delete current file_upload</span>
						{{end}}
						<div class="help">Maximum filesize for web upload: 32M.  file upload</div>
						`
	} else if f.Type == FieldTypeStrMemberGroupRelationship {
		var roles []*UserRole
		json.Unmarshal([]byte(f.Options), &roles)
		params["Selecteds"] = strings.Split(f.DefaultValue, "|")

		var rule FieldTypeTextViewRule
		err := json.Unmarshal([]byte(f.RuleJson), &rule)
		if err == nil {
			params["Rule"] = rule
		}

		params["Groups"] = roles
		tplStr = `<div style="float: left">
			<select name="{{$.FieldType.SystemName}}{{if Equal .Rule.AllowMultiple "1"}}[]{{end}}" {{if Equal .Rule.AllowMultiple "1"}}multiple="multiple"{{end}}>
			{{range .Groups}}
				<option value="{{.UserRoleId}}" {{if InSlice (ToStr .UserRoleId) $.Selecteds}}selected="selected"{{end}}>{{.Name}}</option>
			{{end}}
			</select>
			</div>`
	} else if f.Type == FieldTypeStrRelationship {
		keyValues := make([]map[string]interface{}, 0)
		json.Unmarshal([]byte(f.Options), &keyValues)
		params["Selecteds"] = strings.Split(f.DefaultValue, "|")

		var rule FieldTypeTextViewRule
		err := json.Unmarshal([]byte(f.RuleJson), &rule)
		if err == nil {
			params["Rule"] = rule
			kvs := strings.Split(f.Other, ",")
			if len(kvs) == 2 {
				params["ID"] = kvs[0]
				params["FieldName"] = kvs[1]
			}
			
		}

		params["Maps"] = keyValues
		tplStr = `<div style="float: left">
			<select name="{{$.FieldType.SystemName}}{{if Equal .Rule.AllowMultiple "1"}}[]{{end}}" {{if Equal .Rule.AllowMultiple "1"}}multiple="multiple"{{end}}>
			<option value="">不选</option>
			{{range .Maps}}
				<option value="{{MapGet . $.ID}}" {{if InSlice (ToStr (MapGet . $.ID)) $.Selecteds}}selected="selected"{{end}}>{{MapGet . $.FieldName}}</option>
			{{end}}
			</select>
			</div>`
	}

	strs := utils.Template(header+tplStr+footer, params)
	return strs
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
	ContentHits int64
}

type Language struct {
	LanguageId int64 `orm:"pk;auto"`
	Name string
}

