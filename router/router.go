package router

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/controllers"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/database"
)

var server = controllers.Server{}

func Run() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize(os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	database.Load(server.DB)
	server.Run(":8080")

}