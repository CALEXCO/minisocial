package main

import (
	"log"

	"github.com/CALEXCO/minisocial/internal/config"
	"github.com/CALEXCO/minisocial/internal/platform/database"
	"github.com/CALEXCO/minisocial/internal/server"
)

func main() {

	// First we load the .env which has to be in the main dir
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.Init(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// userStore := user.NewPostgresStore(db)
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	server.Init()
}
