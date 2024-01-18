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

func RemoveEntry(context *gin.Context) {
	id := context.Param("id")

	var input model.Entry

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := input.Remove(id)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusBadRequest, gin.H{"data": "Successfully removed"})
}

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
