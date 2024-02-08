package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/initializers"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/models"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

// SignUpUser    godoc
// @Summary      Register a new user
// @Description	 Create a new user with name, email and password
// @Tags         User
// @Produce      json
// @Param        user  body      models.SignUpInput  true  "Name, Email and Password"
// @Success      200   {object}  models.NewUserResponse
// @Router       /api/v1/auth/register [post]
func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var payload *models.SignUpInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	now := time.Now()
	newUser := models.User{
		Email:     strings.ToLower(payload.Email),
		Password:  hashedPassword,
		Role:      "User",
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		ctx.JSON(http.StatusConflict, gin.H{"message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": "Something bad happened"})
		return
	}

	userResponse := &models.NewUserResponse{
		ID:        newUser.ID,
		Email:     newUser.Email,
		Role:      newUser.Role,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
	}
	ctx.JSON(http.StatusCreated, gin.H{"data": gin.H{"user": userResponse}})
}

// SignInUser    godoc
// @Summary      Login
// @Description	 Login with email and password
// @Tags         User
// @Produce      json
// @Param        User  body      models.SignInInput  true  "Email and Password"
// @Success      200   {object} models.UserResponse
// @Router       /api/v1/auth/login [post]
func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var payload *models.SignInInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or Password"})
		return
	}

	if !user.Status {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Disabled User"})
		return
	}

	if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email or Password"})
		return
	}

	config, _ := initializers.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.SetCookie("Logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("Access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("Refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"Access_token": access_token})
}

// RefreshToken  godoc
// @Summary      RefreshToken
// @Description	 RefreshToken
// @Tags         User
// @Produce      json
// @Success      200
// @Router       /api/v1/auth/refreshToken [GET]
func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "Could not refresh access token"

	cookie, err := ctx.Cookie("Refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": message})
		return
	}

	config, _ := initializers.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	var user models.User
	result := ac.DB.First(&user, "id = ?", fmt.Sprint(sub))
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "The user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
		return
	}

	ctx.SetCookie("Access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("Logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"Access_token": access_token})
}

// LogoutUser    godoc
// @Summary      Logout
// @Description	 Logout
// @Tags         User
// @Produce      json
// @Success      200
// @Router       /api/v1/auth/logout [GET]
func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("Access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("Refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("Logged_in", "", -1, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout Success"})
}
