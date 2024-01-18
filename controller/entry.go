package controller

import (
	"net/http"
	"os"

	"GO-GIN-AIR-POSTGRESQL-DOCKER/helper"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/model"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
)

// CreateEntry             godoc
// @Summary      Create a new entry
// @Description  Responds with the created entry as JSON.
// @Tags         createEntry
// @Accept json
// @Produce      json
// @Router       /entry [post]
// @Param   content body model.EntryInput true "Content of the entry"
// @Success      200  {object}  model.Entry
// @Security ApiKeyAuth
func CreateEntry(context *gin.Context) {
	var input model.Entry
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = user.ID

	savedEntry, err := input.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"data": savedEntry})
}

// UpdateEntry             godoc
// @Summary      Update a new entry
// @Description  Responds with the updated entry as JSON.
// @Tags         updateEntry
// @Accept json
// @Produce      json
// @Router       /entry/{id} [put]
// @Param   id path int true "ID of the entry to be updated"
// @Param   content body model.EntryInput true "Content of the entry"
// @Success      200  {object}  model.Entry
// @Security ApiKeyAuth
func UpdateEntry(context *gin.Context) {
	id := context.Param("id")

	var input model.Entry

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	entry, err := input.Update(id)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusBadRequest, gin.H{"data": entry})

}

// RemoveEntry             godoc
// @Summary      Remove a entry
// @Description  Responds with the success or fail.
// @Tags         removeEntry
// @Accept json
// @Produce      json
// @Router       /entry/{id} [delete]
// @Param   id path int true "ID of the entry to be removed"
// @Success 200 {string} string "success":true,"message":"Entry removed successfully"
// @Security ApiKeyAuth
func RemoveEntry(context *gin.Context) {
	id := context.Param("id")

	var input model.Entry

	err := input.Remove(id)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusBadRequest, gin.H{"data": "Successfully removed"})
}

// GetAllEntries             godoc
// @Summary      Get entries array
// @Description  Responds with the list of all entries as JSON.
// @Tags         GetAllEntries
// @Produce      json
// @Router       /entry [get]
// @Success      200  {array}  model.Entry
// @Security ApiKeyAuth
func GetAllEntries(context *gin.Context) {
	user, err := helper.CurrentUser(context)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"data": user.Entries})
}

func UploadFile(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Failed to upload",
		})
	}

	cld, _ := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("API_KEY"), os.Getenv("API_SECRET"))

	result, err := cld.Upload.Upload(context, file, uploader.UploadParams{
		PublicID: file.Filename,
	})

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Upload to cloudinary failed",
		})
	}

	context.JSON(http.StatusCreated, gin.H{
		"message":   "Successfully uploaded the file",
		"secureURL": result.SecureURL,
		"publicURL": result.URL,
	})
}
