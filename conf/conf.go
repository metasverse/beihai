package conf

import (
	"gopkg.in/yaml.v2"
	"io"
)

type conf struct {
	DB     dbConfig     `yaml:"db"`
	Server serverConfig `yaml:"server"`
	Redis  redisConfig  `yaml:"redis"`
	Jwt    jwtConfig    `yaml:"jwt"`
	Pay    payConfig    `yaml:"pay"`
}

var Instance = conf{}

func Encode(reader io.Reader) error {
	return yaml.NewDecoder(reader).Decode(&Instance)
}
