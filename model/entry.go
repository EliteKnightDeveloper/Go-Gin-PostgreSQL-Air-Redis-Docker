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

type EntryModel struct {
	gorm.Model
	ID uint
}

func (entry *Entry) Save() (*Entry, error) {
	err := database.Database.Create(&entry).Error
	if err != nil {
		return &Entry{}, err
	}
	return entry, nil
}

func (entryModel *EntryModel) Remove() error {
	_, err := FindEntryById(uint(1))
	if err != nil {
		return err
	}
	return err
}

func FindEntryById(entryId uint) (Entry, error) {
	var entry Entry
	err := database.Database.Where("ID=?", entryId).Find(&entry).Error
	if err != nil {
		return Entry{}, err
	}
	return entry, nil
}
