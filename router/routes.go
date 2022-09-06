package router

import (
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"

	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/controllers"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
    r := gin.Default()
    r.Use(func(c *gin.Context) {
        c.Set("db", db)
    })
    r.POST("/users/login", controllers.Login)
    r.POST("/users/register", controllers.CreateUser)
	r.PUT("/users/:userId", controllers.UpdateUser)
	r.DELETE("/users/:userId", controllers.DeleteUser)
    return r
}