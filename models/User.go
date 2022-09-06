package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 			string 		`gorm:"primary_key; unique" json:"id"`
	Username 	string 		`gorm:"size:255;not null;" json:"username"`
	Email 		string 		`gorm:"size:100;not null; unique" json:"email"`
	Password 	string 		`gorm:"size:100;not null;" json:"password"`
	Photos 		[]Photo 	`gorm:"constraint:OnUpdate:CASCADE;" json:"photos"`
	CreatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//Untuk Hash password daru bcrypt
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

//Untuk melakukan validasi password pada saat login
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

//Untuk mengubah password yang akan di save ke dalam bentuk hash password
func (u *User) HashPassword() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

//Untuk inisialisasi data User sebelum di save/update
func (u *User) Initialize() {
	u.ID 		= uuid.New().String()
	u.Username 	= html.EscapeString(strings.TrimSpace(u.Username))
	u.Email 	= html.EscapeString(strings.TrimSpace(u.Email))
	u.Photos	= []Photo{}
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "register":
		if u.ID == "" {
			return errors.New("required user id")
		}else if u.Username == "" {
			return errors.New("required username")
		}else if u.Email == "" {
			return errors.New("required email")
		}else if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}else if u.Password == "" {
			return errors.New("required password")
		}else if len(u.Password) < 6{
			return errors.New("password minimum length 6 characters")
		}

		return nil
	case "update":
		if u.ID == "" {
			return errors.New("required user id")
		}else if u.Username == "" {
			return errors.New("required username")
		}else if u.Email == "" {
			return errors.New("required email")
		}else if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}else if u.Password == "" {
			return errors.New("required password")
		}else if len(u.Password) < 6{
			return errors.New("password minimum length 6 characters")
		}

		return nil
	case "login":
		if u.Password == "" {
			return errors.New("required password")
		}
		if u.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil

	default:
		return nil
	}
}