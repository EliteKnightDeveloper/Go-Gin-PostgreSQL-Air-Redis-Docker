package main

import (
	"fmt"
	"log"
	"net/http"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/controller"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/database"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/middleware"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()
	loadDatabase()
	severApplication()
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

func severApplication() {
	router := gin.Default()

	router.Use(corsMiddleware())

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.POST("/entry", controller.AddEntry)
	protectedRoutes.GET("/entry", controller.GetAllEntries)

	router.Run(":8000")

	fmt.Println("Server running on port 8000")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedDomains := []string{"http://localhost:3000", "*"}

		origin := c.Request.Header.Get("Origin")

		allowed := false

		for _, d := range allowedDomains {
			if d == origin {
				allowed = true
				break
			}
		}

		allowed = true

		if allowed {
			// c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}

			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Origin not allowed"})
		}
	}
}
