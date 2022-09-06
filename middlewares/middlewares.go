package middlewares

import (
	"errors"
	"net/http"

	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/app/auth"
	"github.com/wahyuucandra/task-5-vix-btpns-wahyucandrabuana/app/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, "F", errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
