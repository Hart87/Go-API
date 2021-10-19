package main

import (
	"log"
	"net/http"

	"github.com/hart87/go-api/routes"
)

func main() {
	log.Println("Go-API Running")

	//Routes
	//GET
	http.HandleFunc("/users/all", routes.HandleRoute)
	//GET, PUT, POST, DELETE a User by Id
	http.HandleFunc("/users/", routes.UsersRoute)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}
}
