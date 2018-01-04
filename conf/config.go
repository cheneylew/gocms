package conf

type MYConfig struct {
	ProductName string
	SupportUrl string
	PageLimit int64
}

var GlobalConfig MYConfig

func init() {
	GlobalConfig = MYConfig{
		ProductName:"Hero",
		SupportUrl:"http://www.golangcms.com/",
		PageLimit:3,
	}
}


