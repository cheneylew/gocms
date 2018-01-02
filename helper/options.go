package helper

import (
	"github.com/cheneylew/goutil/utils"
	"strings"
	"github.com/cheneylew/gocms/models"
)

func StringToOptions(str string) []models.Option {
	value := str
	value = utils.Trim(value)
	lines := strings.Split(value, "\n")
	var options []models.Option
	for _, line := range lines {
		line = utils.Trim(line)
		kv := strings.Split(line, "=")
		if len(kv) == 2 {
			key := utils.Trim(kv[0])
			val := utils.Trim(kv[1])
			options = append(options, models.Option{Key:key, Value:val})
		}
	}

	return options
}

func OptionsToString(options []models.Option) string {
	params := utils.TemplateParams()
	params["Options"] = options

	return utils.Template(`{{range .Options}}{{.Key}}={{.Value}}
{{end}}`, params)
}

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