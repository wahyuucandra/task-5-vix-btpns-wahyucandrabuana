package responses

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Membuat record response custom (status, message, dan data)
type Response struct {
	OUT_STAT string 		`json:"OUT_STAT"`
	OUT_MESS string 		`json:"OUT_MESS"`
	OUT_DATA interface{} 	`json:"OUT_DATA"`
}

//Message Response ketika berhasil
func JSON(w http.ResponseWriter, statusCode int, stat string, message string, data interface{}) {
	w.WriteHeader(statusCode)
	
	response := Response{
		OUT_STAT: stat,
		OUT_MESS: string(message),
		OUT_DATA: data,
	}
	
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}

//Message response error
func ERROR(w http.ResponseWriter, statusCode int, stat string, err error) {
	if err != nil {
		JSON(w, statusCode, stat, err.Error(), nil)
		return
	}
	JSON(w, http.StatusBadRequest, "F", "", nil)
}