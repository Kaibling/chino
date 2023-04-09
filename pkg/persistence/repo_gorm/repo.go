package repo_gorm

import (
	"chino/pkg/utils"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&Movie{})
	if err != nil {
		return err
	}
	return nil
}

type dBModel struct {
	ID string `gorm:"primaryKey"`
}

func (db *dBModel) BeforeSave(tx *gorm.DB) (err error) {
	if db.ID == "" {
		id := utils.NewULID()
		db.ID = id.String()
	}
	return nil
}
