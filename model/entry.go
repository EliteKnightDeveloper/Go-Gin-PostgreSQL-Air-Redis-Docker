package model

import (
	"GO-GIN-AIR-POSTGRESQL-DOCKER/database"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	Content string `gorm:"type:text" json:"content"`
	UserID  uint
}

type EntryInput struct {
	Content string `gorm:"type:text" json:"content"`
}

func (entry *Entry) Save() (*Entry, error) {
	err := database.Database.Create(&entry).Error
	if err != nil {
		return &Entry{}, err
	}
	return entry, nil
}

func (entry *Entry) Remove(id string) error {
	var input Entry

	result := database.Database.First(&input, id)
	if result.Error == nil {
		database.Database.Delete(&input, id)
	}

	return result.Error
}

func (entry *Entry) Update(id string) (*Entry, error) {
	var input Entry

	result := database.Database.First(&input, id)
	if result.Error == nil {
		input.Content = entry.Content
		database.Database.Save(&input)
		return &input, nil
	}

	return &Entry{}, result.Error
}
