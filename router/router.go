package router

import (
	"fmt"
	"net/http"
	"pluto-go/controller"

	"github.com/gorilla/mux"
)

func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request Receved:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func GetRouter() *mux.Router {
	router := mux.NewRouter()

	router.Use(loggingMiddleWare)

	router.HandleFunc("/", controller.ServerLive).Methods("GET")
	router.HandleFunc("/signup", controller.SignUp).Methods("POST")

	return router
}
