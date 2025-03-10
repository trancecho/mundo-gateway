package config

import (
	"errors"
	"github.com/spf13/viper"
)

var GlobalConfig *Config

// Config 定义配置结构体
type Config struct {
	Mysql struct {
		Pwd string `mapstructure:"pwd"`
	} `mapstructure:"mysql"`
	configPath string
}

// NewConfig 创建配置实例
func NewConfig() *Config {
	return &Config{}
}

// Init 初始化配置，接收配置文件路径
func (c *Config) Init() error {
	configPath := "./config/config.yaml"
	c.configPath = configPath
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errors.New("配置文件未找到")
		}
		return errors.New("读取配置文件出错: " + err.Error())
	}

	if err := viper.Unmarshal(c); err != nil {
		return errors.New("解析配置文件出错: " + err.Error())
	}
	return nil
}
