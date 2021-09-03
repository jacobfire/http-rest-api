package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jacobfire/http-rest-api/app/apiserver"
	"log"
)

var (
	configPath string
	needMigration bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "config path ")
	flag.BoolVar(&needMigration, "migrate", false, "Start migrations if we need")
}

func main() {
	flag.Parse()

	log.Println("Before initialization")

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	log.Println("Before migrations")
	if needMigration {
		if err = s.Migrate(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Migrations have been started successfully")

		return
	}
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
