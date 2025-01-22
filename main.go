package main

import (
	"fmt"
	"log"
	"net/http"
	"pluto-go/router"
)

func main() {
	fmt.Println("Starting the server")
	fmt.Println("Running the server on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router.GetRouter()))
}
