package configs

import (
	"github.com/BurntSushi/toml"
	"github.com/jacobfire/http-rest-api/app/store/sqlstore"
	"log"
	"sync"
)

var Conf *Config

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store *sqlstore.Config
}

func NewConfig() *Config {
	var once sync.Once
	once.Do(func() {
		Conf = &Config {
			BindAddr: ":8080",
			LogLevel: "info",
			Store:    sqlstore.NewConfig(),
		}
		_, err := toml.DecodeFile("configs/apiserver.toml", Conf)
		if err != nil {
			log.Fatal(err)
		}
	})

	return Conf
}