package config

import (
	"sync"

	"github.com/jinzhu/configor"
)

type DBConfig struct {
	Client     string `default:"mysql" env:"DB_CLIENT"`
	Host       string `default:"127.0.0.1" env:"DB_HOST"`
	Username   string `default:"root" env:"DB_USER"`
	Password   string `required:"true" env:"DB_PASSWORD"`
	Port       uint   `default:"3306" env:"DB_PORT"`
	Database   string `default:"tempo" env:"DB_DATABASE"`
	Migrations struct {
		Path string `default:"./migrations" env:"DB_MIGRATION_PATH"`
	}
	MaxIdleConnections int  `default:"25" env:"DB_MAX_IDLE_CONN"`
	MaxOpenConnections int  `default:"0" env:"DB_MAX_OPEN_CONN"`
	MaxConnLifeTime    int  `default:"90" env:"DB_MAX_CONN_LIFETIME"`
	Debug              bool `default:"false" env:"DB_DEBUG"`
}

type Config struct {
	Service struct {
		Host string `default:"0.0.0.0" env:"SERVICE_HOST"`
		Port string `default:"8080" env:"SERVICE_PORT"`
		Path struct {
			V1 string `default:"/v1" env:"SERVICE_PATH_API"`
		}
	}
	DB       DBConfig
	LogLevel string `default:"INFO" env:"LOG_LEVEL"`
}

var config *Config
var configLock = &sync.Mutex{}

func Instance() Config {
	if config == nil {
		err := Load()
		if err != nil {
			panic(err)
		}
	}
	return *config
}

func Load() error {
	tmpConfig := Config{}
	err := configor.Load(&tmpConfig)
	if err != nil {
		return err
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &tmpConfig

	return nil
}
