package router

import (
	"GO-GIN-AIR-POSTGRESQL-DOCKER/controller"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/middleware"
	"fmt"
	"os"

	_ "GO-GIN-AIR-POSTGRESQL-DOCKER/docs"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SeverApplication() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOrigins = []string{
		"http://localhost:8000",
	}
	config.AddAllowHeaders("Authorization")
	router.Use(cors.New(config))

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicRoutes := router.Group("/api/v1")

	publicRoutes.POST("/register", controller.Register)
	publicRoutes.POST("/login", controller.Login)

	protectedRoutes := router.Group("/api/v1")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())
	protectedRoutes.POST("/entry", controller.CreateEntry)
	protectedRoutes.GET("/entry", controller.GetAllEntries)
	protectedRoutes.PUT("/entry/:id", controller.UpdateEntry)
	protectedRoutes.DELETE("/entry/:id", controller.RemoveEntry)
	protectedRoutes.POST("/file", controller.UploadFile)

	port := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", port)
	router.Run(address)

	fmt.Println("Server running on port 8000")
}
