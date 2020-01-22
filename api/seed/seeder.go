package seed

import (
	"github.com/jinzhu/gorm"
	"go-api/api/models"
	"log"
)

func Load(db *gorm.DB) {

	err := db.Debug().AutoMigrate(&models.User{}, &models.Post{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

}
