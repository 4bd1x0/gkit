package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DebugMode bool                `mapstructure:"debug_mode"`
	Logger    LoggerConfig        `mapstructure:"logger"`
	Web       WebConfig           `mapstructure:"web"`
	DB        map[string]DBConfig `mapstructure:"db"`
	Redis     RedisConfig         `mapstructure:"redis"`
}

type LoggerConfig struct {
	Path    string `mapstructure:"path"`
	Console bool   `mapstructure:"console"`
	Level   string `mapstructure:"level"`
}

type WebConfig struct {
	Host string
	Port int
}

type DBConfig struct {
	Driver   string `mapstructure:"driver"` // mysql, postgres
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	MaxIdle  int    `mapstructure:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open"`
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

var config Config

func IsDebugging() bool {
	return config.DebugMode
}

func GetConfig() *Config {
	return &config
}

func Get(key string) any {
	return viper.Get(key)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func GetStringMap(key string) map[string]any {
	return viper.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}
