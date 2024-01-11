package main

import (
	"log"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/database"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/model"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/router"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	router.SeverApplication()
}

func loadDatabase() {
	database.Connect()
	database.Database.AutoMigrate(&model.User{})
	database.Database.AutoMigrate(&model.Entry{})
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
