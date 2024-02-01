package routes

import (
	"GO-GIN-AIR-POSTGRESQL-DOCKER/controllers"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/middleware"

	"github.com/gin-gonic/gin"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewRoutePostController(postController controllers.PostController) PostRouteController {
	return PostRouteController{postController}
}

func (pc *PostRouteController) PostRoute(rg *gin.RouterGroup) {

	router := rg.Group("posts")
	router.Use(middleware.DeserializeUser())
	router.POST("/", pc.postController.CreatePost)
	router.GET("/", pc.postController.GetPosts)
	router.PUT("/:postId", pc.postController.UpdatePost)
	router.GET("/:postId", pc.postController.GetPostById)
	router.DELETE("/:postId", pc.postController.DeletePost)
}
