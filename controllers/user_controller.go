package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/app/auth"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/helpers/formaterror"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	// Membaca data body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengubah json menjadi object User
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan inisialisasi data user
	user.Initialize()

	//custom error message
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError, "data": nil})
		return
	}

	//Melakukan Hash password
	err = user.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Create(&user)

	//custom error message
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "Success", "data": user})
}

func UpdateUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	// Get model if exist
	var user models.User
	if err := db.Where("id = ?", c.Param("userId")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "User not found", "data": nil})
		return
	}

	// Membaca data body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
	}

	//Mengubah json menjadi object User
	user_input := models.User{}
	user_input.ID = user.ID
	err = json.Unmarshal(body, &user_input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan validasi data user
	err = user_input.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan Hash password
	err = user_input.HashPassword()
	if err != nil {
		log.Fatal(err)
	}

	db.Model(&user).Updates(&user_input)

	//custom error message
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"status": "F", "message": formattedError, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "Success", "data": user_input})
}

func DeleteUser(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", c.Param("userId")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "F", "message": "User not found", "data": nil})
		return
	}

	db.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"status": "T", "message": "Success", "data": nil})
}

//Melakukan login
func Login(c *gin.Context) {
	//Membaca data dari body
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Mengubah json ke objek user
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	//Melakukan perisapan inisialisi dan validasi
	user.Initialize()
	err = user.Validate("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return
	}

	
	//Melakukan pengecekan user di database berdasarkan email
	var user_login models.User
	db := c.MustGet("db").(*gorm.DB)
	if err := db.Where("email = ?", user.Email).First(&user_login).Error; err != nil {
		formattedError := err
		if(err.Error() == "record not found"){
			formattedError = formaterror.ErrorMessage("user not found")
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": formattedError, "data": nil})
		return
	}

	//Melakukan verifikasi password db dan user input
	err = models.VerifyPassword(user_login.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": err.Error(), "data": nil})
		return 
	}

	//Ketika berhasil login akan membuat token jwt
	token, err := auth.CreateToken(user_login.ID)
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": formattedError, "data": nil})
		return
	}
	//Ketika berhasil login memeberikan response success
	c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "F", "message": "success", "data": token})
}
