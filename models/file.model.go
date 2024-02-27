package models

import "mime/multipart"

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}
