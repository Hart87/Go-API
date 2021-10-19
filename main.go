package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hart87/go-api/routes"
)

func main() {
	log.Println(uuid.New())

	//Handlers
	http.HandleFunc("/users/hello", routes.HandleRoute)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
