package database

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/models"
)

var users = []models.User{
	{
		ID: "e771f627-b7ea-4a02-b0c2-3b1c106f1cbc",
		Username: "Steven victor",
		Email:    "steven@gmail.com",
		Password: "password",
	},
	{
		ID: "e771f327-b7ea-4a03-b0c2-3b1c106f1a7a",
		Username: "Martin Luther",
		Email:    "luther@gmail.com",
		Password: "password",
	},
}

var photos = []models.Photo{
	{
		Title:   "Title 1",
		Caption: "Caption 1",
		PhotoUrl: "https://cdn.pixabay.com/photo/2016/10/26/19/00/domain-names-1772240_960_720.png",
	},
	{
		Title:   "Title 2",
		Caption: "Caption 2",
		PhotoUrl: "https://cdn.pixabay.com/photo/2016/10/26/19/00/domain-names-1772240_960_720.png",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Photo{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Photo{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Photo{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		photos[i].UserId = users[i].ID

		err = db.Debug().Model(&models.Photo{}).Create(&photos[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}