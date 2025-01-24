package controller

import (
	"encoding/json"
	"net/http"
	"pluto-go/authentication"
	"pluto-go/database"
	"pluto-go/models"
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

// 1. first write the headers type
// 2. then unmarshal the data coming in thru the request
// 3. and dont allow unkown feild types
// 4. then feed this to the database side
// 5. Get the response as error or ok
// 6. If ok then give them a jwtoken as well
// 7. else throw an error
func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var user models.User

	var decoder *json.Decoder = json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&user); checkError(w, err) {
		return
	}

	// if else logic from database
	err := database.SignUp(&user)

	if checkError(w, err) {
		return
	}

	// create a new token for them
	jwtToken, err := authentication.CreateToken(user.ID.Hex())
	if checkError(w, err) {
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"message": "Signup successful",
			"jwtoken": jwtToken,
		},
	)
}

func checkError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "error",
			"error":   err.Error(),
		})
		return true
	}

	return false
}
