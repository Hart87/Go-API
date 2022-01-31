package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hart87/go-api/handlers"
)

func main() {
	log.Println("Go-API Running")

	//POST. No Auth
	http.HandleFunc("/v1/login", handlers.LoginRoute)

	//POST. No Auth
	http.HandleFunc("/v1/users/new", handlers.NewUserRoute)

	//GET
	http.Handle("/v1/users/all", handlers.IsAuthorized(handlers.GetAll))

	//GET, PUT, DELETE a User by Id
	http.Handle("/v1/users/", handlers.IsAuthorized(handlers.UsersRoute))

	s := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
