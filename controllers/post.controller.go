package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(DB *gorm.DB) PostController {
	return PostController{DB}
}

// CreatePost    godoc
// @Summary      Create a new post
// @Description	 Create a new post with title and content
// @Tags         Post
// @Produce      json
// @Param        post  body      models.CreatePostRequest  true  "Title, Content"
// @Success      200   {object}  models.Post
// @Router       /api/v1/posts [post]
func (pc *PostController) CreatePost(ctx *gin.Context) {
	User := ctx.MustGet("User").(models.User)

	var payload *models.CreatePostRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	now := time.Now()
	newPost := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		User:      User.ID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := pc.DB.Create(&newPost)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Post with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newPost})
}

// GetPosts    	 godoc
// @Summary      Get post lists with user info
// @Description	 Get post lists with user info
// @Tags         Post
// @Produce      json
// @Success      200  {array}  models.Post
// @Router       /api/v1/posts [get]
func (pc *PostController) GetPosts(ctx *gin.Context) {
	User := ctx.MustGet("User").(models.User)

	size, _ := strconv.Atoi(ctx.Query("size"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	var postList []models.PostList
	results := pc.DB.Table("posts").Select("posts.id, posts.title, posts.content,posts.updated_at, users.email, users.id as user").Joins("left join users on posts.user = users.id").Where("posts.user = ?", User.ID).Offset((page - 1) * size).Limit(size).Order("updated_at").Scan(&postList)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(postList), "data": postList})
}

// GetPosts    	 godoc
// @Summary      Get all post lists with user info
// @Description	 Get all post lists with user info
// @Tags         Post
// @Produce      json
// @Success      200  {array}  models.Post
// @Router       /api/v1/posts/all [get]
func (pc *PostController) GetAllPosts(ctx *gin.Context) {
	size, _ := strconv.Atoi(ctx.Query("size"))
	page, _ := strconv.Atoi(ctx.Query("page"))

	var postList []models.PostList
	results := pc.DB.Table("posts").Select("posts.id, posts.title, posts.content,posts.updated_at, users.email, users.id as user").Joins("left join users on posts.user = users.id").Offset((page - 1) * size).Limit(size).Order("updated_at").Scan(&postList)

	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(postList), "data": postList})
}

// UpdatePost    godoc
// @Summary      Update a post
// @Description	 Update a post
// @Tags         Post
// @Produce      json
// @Param        post  body      models.UpdatePost  true  "Title, Content"
// @Param   		 postId path string true "ID of the entry to be updated"
// @Success      200   {object}  models.Post
// @Router       /api/v1/posts/{postId} [put]
func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	User := ctx.MustGet("User").(models.User)

	var payload *models.UpdatePost
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedPost models.Post
	result := pc.DB.First(&updatedPost, "id = ?", postId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}
	now := time.Now()
	postToUpdate := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		User:      User.ID,
		CreatedAt: updatedPost.CreatedAt,
		UpdatedAt: now,
	}

	pc.DB.Model(&updatedPost).Updates(postToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedPost})
}

// GetPost    	 godoc
// @Summary      Get a post
// @Description	 Get a post
// @Tags         Post
// @Produce      json
// @Param   		 postId path string true "ID of the entry to be retrived"
// @Success      200
// @Router       /api/v1/posts/{postId} [get]
func (pc *PostController) GetPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	var post models.Post
	result := pc.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": post})
}

// DeletePost    godoc
// @Summary      Delete a post
// @Description	 Delete a post
// @Tags         Post
// @Produce      json
// @Param   	 	 postId path string true "ID of the entry to be deleted"
// @Success      200
// @Router       /api/v1/posts/{postId} [delete]
func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	result := pc.DB.Delete(&models.Post{}, "id = ?", postId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
