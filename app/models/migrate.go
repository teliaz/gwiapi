package models

import (
	"github.com/jinzhu/gorm"
)

func DBMigrate(db *gorm.DB) *gorm.DB {

	db.DropTableIfExists(&Asset{}, &Audience{}, &Insight{}, &Chart{}, &User{}, &Data{})

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Asset{})
	db.AutoMigrate(&Audience{})
	db.AutoMigrate(&Insight{})
	db.AutoMigrate(&Chart{})

	db.AutoMigrate(&Data{})
	// db.Model(&Chart{}).AddForeignKey("AssetId", "Assets(id)", "CASCADE", "CASCADE")
	return db
}
