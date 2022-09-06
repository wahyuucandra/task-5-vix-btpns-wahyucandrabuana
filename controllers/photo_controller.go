package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/app/responses"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/models"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/helpers/formaterror"

)

func (server *Server) CreatePhoto(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, "F", err)
		return
	}
	user.Prepare()
	// err = user.Validate("")
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnprocessableEntity, err)
	// 	return
	// }
	userCreated, err := user.Register(server.DB)

	if err != nil {

		formattedError := formaterror.ErrorMessage(err.Error())

		responses.ERROR(w, http.StatusInternalServerError, "F", formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%s", r.Host, r.RequestURI, userCreated.ID))
	responses.JSON(w, http.StatusCreated, "T", "success", userCreated)
}

func (server *Server) GetPhotos(w http.ResponseWriter, r *http.Request) {

	photo := models.Photo{}

	users, err := photo.GetPhotos(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, "F", err)
		return
	}
	responses.JSON(w, http.StatusOK, "T", "success", users)
}