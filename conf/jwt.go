package conf

type jwtConfig struct {
	SecretKey string `yaml:"secret_key"`
	Expire    int64  `yaml:"expire"`
}
