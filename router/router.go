package router

import (
	"GO-GIN-AIR-POSTGRESQL-DOCKER/controller"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/middleware"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func SeverApplication() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	publicRoutes := router.Group("/auth")

	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.POST("/entry", controller.AddEntry)
	protectedRoutes.GET("/entry", controller.GetAllEntries)
	protectedRoutes.POST("/file", controller.UploadFile)

	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)
	fmt.Println(address)
	router.Run(address)

	fmt.Println("Server running on port 8000")
}
