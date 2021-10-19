package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hart87/go-api/handlers"
)

func main() {
	log.Println("Go-API Running")

	//GET
	http.HandleFunc("/users/all", handlers.HandleRoute)
	//GET, PUT, POST, DELETE a User by Id
	http.HandleFunc("/users/", handlers.UsersRoute)

	s := &http.Server{
		Addr:         ":8082",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	log.Fatal(s.ListenAndServe())

}
