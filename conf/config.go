package conf

type MYConfig struct {
	ProductName string
	SupportUrl string

}

var GlobalConfig MYConfig

func init() {
	GlobalConfig = MYConfig{
		ProductName:"Hero",
		SupportUrl:"http://www.baidu.com/",
	}
}


