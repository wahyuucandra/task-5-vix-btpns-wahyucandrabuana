package models

import (
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 			int 		`gorm:"primary_key;auto_increment" json:"id"`
	Username 	string 		`gorm:"size:255;not null;unique" json:"username"`
	Email 		string 		`gorm:"size:100;not null;unique" json:"email"`
	Password 	string 		`gorm:"size:100;not null;" json:"password"`
	Photos 		[]Photo 	`gorm:"constraint:OnUpdate:CASCADE;" json:"photos"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID 		= 0
	u.Username 	= html.EscapeString(strings.TrimSpace(u.Username))
	u.Email 	= html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Register(db *gorm.DB) (*User, error) {

	var err error = nil
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}