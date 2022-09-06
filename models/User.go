package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
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
	id 	:= uuid.New()
	u.ID 		= id.String()
	u.Username 	= html.EscapeString(strings.TrimSpace(u.Username))
	u.Email 	= html.EscapeString(strings.TrimSpace(u.Email))
	u.Photos	= []Photo{}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
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

func (u *User) Register(db *gorm.DB) (*User, error) {

	var err error = nil
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	if len(users) > 0 {
		for i := range users {
			err := db.Debug().Model(&Photo{}).Where("user_id = ?", users[i].ID).Take(&users[i].Photos).Error
			if err != nil {
				return &[]User{}, err
			}
		}
	}
	return &users, err
}

func (u *User) UpdateUser(db *gorm.DB, id string) (*User, error) {

	// Sebelum melakukan update password dilakukan hash terlebih dahulu
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	//Melakukan sql (gorm) update data user
	db = db.Debug().Model(&User{}).Where("id = ?", id).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password"	: u.Password,
			"username"	: u.Username,
			"email"		: u.Email,
			"update_at"	: time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	//Melakukan update user dengan data yang baru
	err = db.Debug().Model(&User{}).Where("id = ?", id).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	//Menghapus data user
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
