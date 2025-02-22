package setting

import (
	"gopkg.in/ini.v1"
)

var Conf = new(AppConfig)

// AppConfig 应用程序配置
type AppConfig struct {
	Release   bool   `ini:"release"`
	Port      int    `ini:"port"`
	SecretId  string `ini:"secretId"`
	SecretKey string `ini:"secretKey"`
}

func Init(file string) error {
	return ini.MapTo(Conf, file)
}

// GetSecretId 获取SecretId
func (c *AppConfig) GetSecretId() string {
	return c.SecretId
}

// GetSecretKey 获取SecretKey
func (c *AppConfig) GetSecretKey() string {
	return c.SecretKey
}
