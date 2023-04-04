package config

import (
	"landlord/program/util"

	"github.com/BurntSushi/toml"
	"github.com/google/logger"
)

//订制配置文件解析载体
type Config struct {
	Database *Database
}

type Database struct {
	Host     string
	Port     int
	DbName   string
	Username string
	Password string
}

var Con *Config = new(Config)

func init() {
	//读取配置文件
	_, err := toml.DecodeFile(util.GetConfigFilePath(), Con)
	if err != nil {
		logger.Info(err)
	}

}
