package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Web      WebConfig      `mapstructure:"web"`
	Database DatabaseConfig `mapstructure:"database"`
	Midtrans MidtransConfig `mapstructure:"midtrans"`
}

type MidtransConfig struct {
	ServerKey   string `mapstructure:"server_key"`
	ClientKey   string `mapstructure:"client_key"`
	Environment string `mapstructure:"environment"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
}

type WebConfig struct {
	Port        string `mapstructure:"port"`
	Prefork     bool   `mapstructure:"prefork"`
	AppPassword string `mapstructure:"app_password"`
}

type DatabaseConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Name     string `mapstructure:"name"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.Username, d.Password, d.Host, d.Port, d.Name)
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
