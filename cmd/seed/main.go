package main

import (
	"log"
	"template-go-api/config"
	"template-go-api/database/seeders"
)

func main() {
	config.LoadEnv()

	db := config.ConnectDB()

	log.Println("Running seeders...")
	seeders.RunAll(db)
	log.Println("Seeding complete.")
}
