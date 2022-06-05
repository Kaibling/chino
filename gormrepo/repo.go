package gormrepo

import (
	"chino/lib/utils"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Movie{})
	if err != nil {
		return err
	}
	return nil
}

type DBModel struct {
	ID string `gorm:"primaryKey"`
}

func (db *DBModel) BeforeSave(tx *gorm.DB) (err error) {
	if db.ID == "" {
		id := utils.NewULID()
		db.ID = id.String()
	}
	return nil
}
