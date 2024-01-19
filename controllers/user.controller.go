package controllers

import (
	"net/http"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

// GetCurrentUserInfo  godoc
// @Summary      GetCurrentUserInfo
// @Description	 GetCurrentUserInfo
// @Tags         User
// @Produce      json
// @Success      200 {object}  models.NewUserResponse
// @Router       /api/v1/users/me [GET]
func (uc *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	userResponse := &models.NewUserResponse{
		ID:        currentUser.ID,
		Name:      currentUser.Name,
		Email:     currentUser.Email,
		Role:      currentUser.Role,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}
