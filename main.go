package main

import (
	"log"
	"net/http"

	"github.com/hart87/go-api/routes"
)

func main() {
	log.Println("Go-API Running")

	//Routes
	http.HandleFunc("/users/hello", routes.HandleRoute)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
