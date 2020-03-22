package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"./app"
	"./db"
)

func main() {
	database, err := db.InitDbConnection()
	if err != nil {
		log.Fatal("Database connection failed: %s", err.Error())
	}

	app := &app.App{
		Router:   mux.NewRouter().StrictSlash(true),
		Database: database,
	}

	app.SetupRouter()

	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
