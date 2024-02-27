package controllers

import (
	"net/http"
	"time"

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
	User := ctx.MustGet("User").(models.User)

	userResponse := &models.NewUserResponse{
		ID:        User.ID,
		Email:     User.Email,
		Role:      User.Role,
		CreatedAt: User.CreatedAt,
		UpdatedAt: User.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{"data": gin.H{"user": userResponse}})
}

// ApproveUser	 godoc
// @Summary      ApproveUser
// @Description	 ApproveUser
// @Tags         User
// @Produce      json
// @Param   	 userId path string true "ID of the user to be approved"
// @Success      200
// @Router       /api/v1/users/approve/{userId} [POST]
func (uc *UserController) ApproveUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var user models.User
	result := uc.DB.First(&user, "id=?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No user with this id"})
		return
	}

	now := time.Now()
	userToUpdate := models.User{
		UpdatedAt: now,
		Status:    true,
	}

	uc.DB.Model(&user).Updates(userToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"message": "User Approve Success"})
}

// DisableUser	 godoc
// @Summary      DisableUser
// @Description	 DisableUser
// @Tags         User
// @Produce      json
// @Param   	 userId path string true "ID of the user to be disabled"
// @Success      200
// @Router       /api/v1/users/disable/{userId} [POST]
func (uc *UserController) DisableUser(ctx *gin.Context) {
	userId := ctx.Param("userId")

	var user models.User
	result := uc.DB.First(&user, "id=?", userId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No user with this id"})
		return
	}

	now := time.Now()
	userToUpdate := models.User{
		UpdatedAt: now,
		Status:    false,
	}

	uc.DB.Model(&user).Updates(userToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"message": "User Disabled Success"})
}
