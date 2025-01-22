package controller

import (
	"encoding/json"
	"net/http"
)

func ServerLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]string{
			"message": "Server Live",
		},
	)
}

func SignUp(w http.ResponseWriter , r *http.Request){
	// we have to first get the details of the user from the body
}