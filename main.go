package main

import (
	"log"
	"net/http"
	"time"

	"github.com/hart87/go-api/handlers"
)

func main() {
	log.Println("Go-API Running")

	//TEST ROUTE
	http.HandleFunc("/test", handlers.TestRoute)
	//GET
	http.HandleFunc("/users/all", handlers.GetAllUsers)
	//GET, PUT, POST, DELETE a User by Id
	http.HandleFunc("/users/", handlers.UsersRoute)

	s := &http.Server{
		Addr:         ":8083",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	log.Fatal(s.ListenAndServe())

}
