package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func SetupDB() *gorm.DB {
	DB_HOST 	:= "127.0.0.1"
	DB_DRIVER 	:= "postgres"
	DB_USER 	:= "postgres"
	DB_PASSWORD := "admin"
	DB_NAME		:= "btpn"
	DB_PORT     := "5432" 

	URL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DB_HOST, DB_PORT, DB_USER, DB_NAME, DB_PASSWORD)
	db, err := gorm.Open("postgres", URL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database", DB_DRIVER)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database", DB_DRIVER)
	}
	return db
}
