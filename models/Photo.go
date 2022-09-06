package models

import (
	"html"
	"strings"

	"github.com/jinzhu/gorm"
)

type Photo struct {
	ID 			int 	`gorm:"primary_key;auto_increment" json:"id"`
	Title 		string 	`gorm:"size:100;not null" json:"title"`
	Caption 	string 	`gorm:"size:255;not null" json:"caption"`
	PhotoUrl 	string 	`gorm:"size:255;not null;" json:"photo_url"`
	UserId 		string 	`gorm:"not null" json:"user_id"`
	Author      User   	`gorm:"author"`
}

func (p *Photo) Prepare() {
	p.ID 		= 0
	p.Title 	= html.EscapeString(strings.TrimSpace(p.Title))
	p.Caption 	= html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl 	= html.EscapeString(strings.TrimSpace(p.PhotoUrl))
}

func (u *Photo) GetPhotos(db *gorm.DB) (*[]Photo, error) {
	var err error
	photos := []Photo{}
	err = db.Debug().Model(&Photo{}).Limit(100).Find(&photos).Error
	if err != nil {
		return &[]Photo{}, err
	}
	if len(photos) > 0 {
		for i := range photos {
			err := db.Debug().Model(&User{}).Where("id = ?", photos[i].UserId).Take(&photos[i].Author).Error
			if err != nil {
				return &[]Photo{}, err
			}
		}
	}
	return &photos, err
}