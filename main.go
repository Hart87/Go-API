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
	http.HandleFunc("/v1/users/all", handlers.GetAllUsers)
	//GET, PUT, POST, DELETE a User by Id
	http.HandleFunc("/v1/users/", handlers.UsersRoute)

	s := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	log.Fatal(s.ListenAndServe())

}
