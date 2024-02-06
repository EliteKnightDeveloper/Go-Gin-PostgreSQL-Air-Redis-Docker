package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Title     string    `gorm:"uniqueIndex;not null" json:"title,omitempty"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	User      uuid.UUID `gorm:"type:uuid not null" json:"user,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at,omitempty"`
}

type CreatePostRequest struct {
	Title   string `json:"title"  binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdatePost struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

type PostList struct {
	ID      uuid.UUID `json:"ID"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
	Email   string    `json:"email,omitempty"`
	User    string    `json:"user,omitempty"`
}
