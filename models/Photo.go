package models

import (
	"html"
	"strings"
)

type Photo struct {
	ID 			int 	`gorm:"primary_key;auto_increment" json:"id"`
	Title 		string 	`gorm:"size:255;not null;unique" json:"username"`
	Caption 	string 	`gorm:"size:100;not null;unique" json:"email"`
	PhotoUrl 	string 	`gorm:"size:100;not null;" json:"password"`
	UserId 		string 	`gorm:"not null" json:"user_id"`
	Author      User   	`gorm:"author"`
}

func (p *Photo) Prepare() {
	p.ID 		= 0
	p.Title 	= html.EscapeString(strings.TrimSpace(p.Title))
	p.Caption 	= html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoUrl 	= html.EscapeString(strings.TrimSpace(p.PhotoUrl))
	p.UserId 	= html.EscapeString(strings.TrimSpace(p.UserId))
	p.Author 	= User{}
}