package helper

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func MakeRespon(w http.ResponseWriter, status int, msg string, data interface{}) {
	w.Header().Set("Content-type", "application/json")
	var response Response
	response.Status = status
	response.Message = msg
	response.Data = data
	userJson, _ := json.Marshal(response)
	w.WriteHeader(status)
	w.Write(userJson)
}
