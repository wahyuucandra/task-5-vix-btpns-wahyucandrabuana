package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/models"
	"golang.org/x/crypto/bcrypt"
)
type CreateUserInput struct {
    Username 	string 		`json:"username"`
	Email 		string 		`json:"email"`
	Password 	string 		`json:"password"`
}

type UpdateUserInput struct {
    Username 	string 		`json:"username"`
	Email 		string 		`json:"email"`
	Password 	string 		`json:"password"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CreateUser(c *gin.Context) {
	// Validate input
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	id 	:= uuid.New()
	hashedPassword, err := Hash(input.Password)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
	input.Password = string(hashedPassword)

	// Create user
	user := models.User{
		ID: id.String(),
		Username: input.Username,
		Email: input.Email,
		Password: input.Password,
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Create(&user)

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "Success", "data": user})
}

func UpdateUser(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    // Get model if exist
    var user models.User
    if err := db.Where("id = ?", c.Param("userId")).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak ditemukan!!"})
        return
    }

    // Validate input
    var input UpdateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	hashedPassword, err := Hash(input.Password)
	if err != nil {
		log.Fatal("This is the error:", err)
	}
	input.Password = string(hashedPassword)

    var updatedInput models.User
    updatedInput.Email 		= input.Email
    updatedInput.Username 	= input.Username
    updatedInput.Password 	= input.Password

    db.Model(&user).Updates(updatedInput)

    c.JSON(http.StatusOK, gin.H{"status": "T", "message": "Success", "data": user})
}

func DeleteUser(c *gin.Context) {

    db := c.MustGet("db").(*gorm.DB)
    var user models.User
    if err := db.Where("id = ?", c.Param("userId")).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak ditemukan!!"})
        return
    }

    db.Delete(&user)

    c.JSON(http.StatusOK, gin.H{"status": "T", "message": "Success", "data": nil})
}



