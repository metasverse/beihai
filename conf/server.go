package conf

type serverConfig struct {
	Addr      string `yaml:"addr"`
	Domain    string `yaml:"domain"`
	SecretKey string `yaml:"secret_key"`
}
