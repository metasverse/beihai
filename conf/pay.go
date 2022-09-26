package conf

type payConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
	MID       string `yaml:"mid"`
	TID       string `yaml:"tid"`
	PayURL    string `yaml:"pay_url"`
	QueryURL  string `yaml:"query_url"`
}
