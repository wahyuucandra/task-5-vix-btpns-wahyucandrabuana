package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/app/auth"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/app/responses"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/helpers/formaterror"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/models"
	"golang.org/x/crypto/bcrypt"
)

//Fungsi Register User
func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {

	//membaca response body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
	}

	//Mengubah json menjadi object User
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}

	//Melakukan prepare inisialisi data
	user.Prepare()
	//Melakukan pengecekan validasi ketika register
	err = user.Validate("register")
	
	//Menampilkan error dari validasi
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}
	userCreated, err := user.Register(server.DB)

	//custom error message
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, "F", formattedError)
		return
	}

	//Memberikan respose ketika berhasil
	w.Header().Set("Location", fmt.Sprintf("%s%s/%s", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, "T", "user successfully registered", userCreated)
}

//Menampilkan semua user yang terdaftar ke database
func (server *Server) GetUsers(w http.ResponseWriter, r *http.Request) {

	user := models.User{}
	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, "F", err)
		return
	}
	responses.JSON(w, http.StatusOK, "T", "Success", users)
}

//Melakukan pengecekan login dan membuat token jwt
func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	//Melakukan pengecekan user di database berdasarkan email
	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		var formattedError error
		if(err.Error() == "record not found"){
			formattedError = formaterror.ErrorMessage("user not found")
		}
		return "", formattedError
	}

	//Melakukan verifikasi password db dan user input
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	//Ketika berhasil login akan membuat token jwt
	return auth.CreateToken(user.ID)
}


//Melakukan login
func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	//Membaca data dari bosy
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}

	//Mengubah json ke objek user
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}

	//Melakukan perisapan inisialisi dan validasi
	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}

	//Melakukan pengecekan login
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.ErrorMessage(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", formattedError)
		return
	}
	//Ketika berhasil login memeberikan response success
	responses.JSON(w, http.StatusOK, "T", "login successfully", 
		struct {Token string `json:"token"`}{ Token: token,})
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, "F", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, "F", errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		responses.ERROR(w, http.StatusUnauthorized, "F", errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}
	// updatedUser, err := user.UpdateUser(server.DB, uid)
	// if err != nil {
	// 	formattedError := formaterror.ErrorMessage()(err.Error())
	// 	responses.ERROR(w, http.StatusInternalServerError, formattedError)
	// 	return
	// }
	// responses.JSON(w, http.StatusOK, updatedUser)
}

// func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {

// 	vars := mux.Vars(r)

// 	user := models.User{}

// 	uid, err := strconv.ParseUint(vars["id"], 10, 32)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusBadRequest, err)
// 		return
// 	}
// 	tokenID, err := auth.ExtractTokenID(r)
// 	if err != nil {
// 		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
// 		return
// 	}
// 	if tokenID != 0 && tokenID != uint32(uid) {
// 		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
// 		return
// 	}
// 	_, err = user.DeleteAUser(server.DB, uint32(uid))
// 	if err != nil {
// 		responses.ERROR(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
// 	responses.JSON(w, http.StatusNoContent, "")
// }

