package main

import (
	"log"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/database"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/model"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/router"

	"github.com/joho/godotenv"
)

// @title           Go-Gin-PostgreSQL-Air
// @version         1.0
// @description     A CRUD boilerplate project in GO, Gin, PostgreSQL, Air framework.

// @contact.name   Nguyen Hieu
// @contact.email  trungquann411@gmail.com

// @host      localhost:8000
// @BasePath  /api/v1
// @securitydefinitions.apikey ApiKeyAuth
// @in       Header
// @name      Authorization
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
