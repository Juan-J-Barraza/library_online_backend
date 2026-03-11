package main

import (
	"libraryOnline/config/db"
	"libraryOnline/routers"
	"libraryOnline/utils"
	"log"
)

func main() {

	dataBase, err := db.SetDatabase()
	if err != nil {
		log.Fatalf("error to init database %v", err)
	}

	db.Seed(dataBase)
	db.SeedLoans(dataBase)

	app, err := utils.InitFiber()
	if err != nil {
		log.Fatalf("error to init app fiber %v", err)
	}

	routers.SetRouters(dataBase, app)

	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("error to init server %v", err)
	}
}
