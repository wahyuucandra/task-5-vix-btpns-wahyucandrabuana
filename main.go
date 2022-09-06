package main

import (
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/database"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/models"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/router"
)

func main() {
	db := database.SetupDB()
    db.AutoMigrate(&models.User{})

    r := router.SetupRoutes(db)
    r.Run()
}