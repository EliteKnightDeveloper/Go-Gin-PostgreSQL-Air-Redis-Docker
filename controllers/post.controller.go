package controllers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/models"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostController struct {
	DB *gorm.DB
}

func NewPostController(DB *gorm.DB) PostController {
	return PostController{DB}
}

// CreatePost godoc
// @Summary      Create a new post
// @Description  Create a new post with title and content
// @Tags         Post
// @Produce      json
// @Param        post  body models.CreatePost true "Title, Content"
// @Param        file formData file true "File"
// @Success      200   {object}  models.Post
// @Router       /api/v1/posts [post]
func (pc *PostController) CreatePost(ctx *gin.Context) {
	user := ctx.MustGet("User").(models.User)

	file, _, err := ctx.Request.FormFile("file")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is provided",
		})
		return
	}
	uploadUrl, err := services.UploadHelper(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "File upload is failed."})
		return
	}

	var payload *models.CreatePost
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	now := time.Now()
	newPost := models.Post{
		Title:     payload.Content,
		Content:   payload.Title,
		User:      user.ID,
		CreatedAt: now,
		UpdatedAt: now,
		File:      uploadUrl,
	}

	result := pc.DB.Create(&newPost)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"message": "Post with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newPost})
}

// GetPosts    	 godoc
// @Summary      Get post lists with user info
// @Description	 Get post lists with user info
// @Tags         Post
// @Produce      json
// @Success      200  {array}  models.Post
// @Param   	 size query string true "Page size"
// @Param   	 page query string true "Page"
// @Router       /api/v1/posts [get]
func (pc *PostController) GetPosts(ctx *gin.Context) {
	User := ctx.MustGet("User").(models.User)

	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Page size is not provided."})
		return
	}
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Page number is not provided."})
		return
	}

	var postList []models.PostList
	results := pc.DB.Table("posts").Select("posts.id, posts.title, posts.content,posts.updated_at, users.email, users.id as user").Joins("left join users on posts.user = users.id").Where("posts.user = ?", User.ID).Offset((page - 1) * size).Limit(size).Order("updated_at").Scan(&postList)
	if results.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": results.Error})
		return
	}

	var pageSize = 0

	if len(postList)/size == 0 {
		pageSize = 1
	} else {
		pageSize = len(postList) / size
	}

	ctx.JSON(http.StatusOK, gin.H{"total": len(postList), "page": pageSize, "data": postList})
}

// GetPosts    	 godoc
// @Summary      Get all post lists with user info
// @Description	 Get all post lists with user info
// @Tags         Post
// @Produce      json
// @Success      200  {array}  models.Post
// @Param   	 size query int true "Page size"
// @Param   	 page query int true "Page"
// @Router       /api/v1/posts/all [get]
func (pc *PostController) GetAllPosts(ctx *gin.Context) {
	size, err := strconv.Atoi(ctx.Query("size"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Page size is not provided."})
		return
	}
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Page number is not provided."})
		return
	}

	var postList []models.PostList
	results := pc.DB.Table("posts").Select("posts.id, posts.title, posts.content,posts.updated_at, users.email, users.id as user").Joins("left join users on posts.user = users.id").Offset((page - 1) * size).Limit(size).Order("updated_at").Scan(&postList)

	if results.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": results.Error})
		return
	}

	var pageSize = 0

	if len(postList)/size == 0 {
		pageSize = 1
	} else {
		pageSize = len(postList) / size
	}

	ctx.JSON(http.StatusOK, gin.H{"total": len(postList), "page": pageSize, "data": postList})
}

// UpdatePost    godoc
// @Summary      Update a post
// @Description	 Update a post
// @Tags         Post
// @Produce      json
// @Param        post body models.UpdatePost true "Title, Content"
// @Param   	 postId path string true "ID of the entry to be updated"
// @Success      200 {object} models.Post
// @Router       /api/v1/posts/{postId} [put]
func (pc *PostController) UpdatePost(ctx *gin.Context) {
	postId := ctx.Param("postId")
	user := ctx.MustGet("User").(models.User)

	var payload *models.UpdatePost
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var updatedPost models.Post
	result := pc.DB.First(&updatedPost, "id = ?", postId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No post with that Id exists"})
		return
	}
	now := time.Now()
	postToUpdate := models.Post{
		Title:     payload.Title,
		Content:   payload.Content,
		User:      user.ID,
		CreatedAt: updatedPost.CreatedAt,
		UpdatedAt: now,
	}

	pc.DB.Model(&updatedPost).Updates(postToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"data": updatedPost})
}

// GetPost    	 godoc
// @Summary      Get a post
// @Description	 Get a post
// @Tags         Post
// @Produce      json
// @Param   	 postId path string true "ID of the entry to be retrived"
// @Success      200
// @Router       /api/v1/posts/{postId} [get]
func (pc *PostController) GetPostById(ctx *gin.Context) {
	postId := ctx.Param("postId")

	var post models.Post
	result := pc.DB.First(&post, "id = ?", postId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No post with that Id exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": post})
}

// DeletePost    godoc
// @Summary      Delete a post
// @Description	 Delete a post
// @Tags         Post
// @Produce      json
// @Param   	 postId path string true "ID of the entry to be deleted"
// @Success      200
// @Router       /api/v1/posts/{postId} [delete]
func (pc *PostController) DeletePost(ctx *gin.Context) {
	postId := ctx.Param("postId")

	result := pc.DB.Delete(&models.Post{}, "id = ?", postId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No post with that Id exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
