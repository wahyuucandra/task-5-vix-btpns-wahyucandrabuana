package router

import (
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"

	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/controllers"
    "github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/middlewares"
)

//Melakukan set up routes end point
func SetupRoutes(db *gorm.DB) *gin.Engine {
    r := gin.Default()
    r.Use(func(c *gin.Context) {
        c.Set("db", db)
    })

    //User
    r.POST("/users/login", controllers.Login)
    r.POST("/users/register", controllers.CreateUser)
	r.PUT("/users/:userId", controllers.UpdateUser)
	r.DELETE("/users/:userId", controllers.DeleteUser)

    //PhotoUrl
    r.GET("/photos", controllers.GetPhoto)
    
    //Photo Url yang membutuhkan token jwt
    secured := r.Group("/").Use(middlewares.Auth())
    {
        secured.POST("/photos", controllers.CreatePhoto)
        secured.PUT("/photos/:photoId", controllers.UpdatePhoto)
        secured.DELETE("/photos/:photoId", controllers.DeletePhoto)
    }

    return r
}