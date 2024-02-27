package services

import (
	"GO-GIN-AIR-POSTGRESQL-DOCKER/models"
	"GO-GIN-AIR-POSTGRESQL-DOCKER/utils"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type mediaUpload interface {
	FileUpload(file models.File) (string, error)
}

type media struct{}

func NewMediaUpload() mediaUpload {
	return &media{}
}

func (*media) FileUpload(file models.File) (string, error) {
	//validate
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, err := utils.UploadHelper(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}
