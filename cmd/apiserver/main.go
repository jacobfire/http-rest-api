package main

import (
	"flag"
	"fmt"
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

	s := apiserver.New()
	log.Println("Before migrations")
	if needMigration {
		if err := s.Migrate(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Migrations have been started successfully")

		return
	}
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
