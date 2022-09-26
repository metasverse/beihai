package conf

import "fmt"

type redisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (r redisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}
