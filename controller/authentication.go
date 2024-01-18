package controller

import (
	"net/http"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/helper"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/model"

	"github.com/gin-gonic/gin"
)

// Register             godoc
// @Summary      Register a new user with username and password.
// @Description  Responds with the created user as JSON.
// @Tags         Register
// @Accept json
// @Produce      json
// @Router       /register [post]
// @Param   content body model.UserInput true "Username and password"
// @Success      200  {object} model.User
func Register(context *gin.Context) {
	var input model.UserInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Username: input.Username,
		Password: input.Password,
	}

	savedUser, err := user.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user": savedUser})
}

// Login             godoc
// @Summary      Login with username and password.
// @Description  Responds with the token string.
// @Tags         Login
// @Accept json
// @Produce      json
// @Router       /login [post]
// @Param   content body model.UserInput true "Username and password"
// @Success      200  {object} map[token]string
func Login(context *gin.Context) {
	var input model.UserInput

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := model.FindUserByUsername(input.Username)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidatePassword(input.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := helper.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": "Bearer " + jwt})
}
