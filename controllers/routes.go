package controllers

import (
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/middlewares"
)

func (s *Server) initializeRoutes() {

	// Login Route
	s.Router.HandleFunc("/users/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users/register", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	

	//Users routes
	s.Router.HandleFunc("/photos", middlewares.SetMiddlewareJSON(s.GetPhotos)).Methods("GET")
	

}