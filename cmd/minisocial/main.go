package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/CALEXCO/minisocial/internal/config"
	"github.com/CALEXCO/minisocial/internal/platform/database"
	"github.com/CALEXCO/minisocial/internal/server"
	"github.com/CALEXCO/minisocial/internal/user"
)

func main() {
	fmt.Println(strings.Contains("alex@gmail.com", "@"))
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

	userStore := user.NewPostgresStore(db)
	userService := user.NewService(userStore)
	userHandler := user.NewUserHandler(userService)
	// userStore := user.NewPostgresStore(db)
	if err := database.Migrate(db); err != nil {
		log.Println(err)
	}

	if err := server.Init(conf, userHandler); err != nil {
		log.Fatal(err)
	}

}
