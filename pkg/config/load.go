package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

type PostgresConfig struct {
	Host     string `yaml:"POSTGRES_HOST"`
	User     string `yaml:"POSTGRES_USER"`
	Password string `yaml:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DB"`
	Port     string `yaml:"POSTGRES_PORT"`
	Timezone string `yaml:"TIMEZONE"`
}

type DiscountConfig struct {
	ExpireTime int `yaml:"expire_time"`
	CodeLength int `yaml:"code_length"`
}

type Config struct {
	ServerPort     int            `yaml:"port"`
	Token          string         `yaml:"token"`
	DiscountConfig DiscountConfig `yaml:"discount"`
	PostgresConfig PostgresConfig `yaml:"postgres"`
}

func LoadConfig(path string) (config *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	config = &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
