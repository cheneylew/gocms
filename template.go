package main

import (
	"html/template"
	"github.com/cheneylew/gocms/models"
	"bytes"
	"github.com/cheneylew/goutil/utils"
)

func TemplateMain()  {
	a := models.FieldTypeCheckBox{}
	t := template.Must(template.New("test").Parse(`<li>
	<label for="default">Default State</label>
	<select name="default">
	{{range .Options}}
	<option value="{{.Value}}" {{if eq .Value $.Option.Value}}selected="selected"{{end}}>{{.Name}}</option>
	{{end}}
	</select>
	</li>`))
	mp := make(map[string]interface{}, 0)
	mp["Options"] = a.Options()
	mp["Option"] = a.DefaultValue()

	buf := bytes.NewBufferString("")
	t.Execute(buf, mp)
	utils.JJKPrintln(buf)
}
