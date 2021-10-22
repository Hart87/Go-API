package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hart87/go-api/db"
	"github.com/hart87/go-api/handlers"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.Println("Go-API Running")

	runMongo()

	//GET
	http.HandleFunc("/users/all", handlers.HandleRoute)
	//GET, PUT, POST, DELETE a User by Id
	http.HandleFunc("/users/", handlers.UsersRoute)

	s := &http.Server{
		Addr:         ":8082",
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	log.Fatal(s.ListenAndServe())

}

func runMongo() {
	clientOptions := options.Client().ApplyURI(
		db.CONNECTION_URI + db.CONNECTION_PORT)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)

	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("connected successfully to MongoDB @ " + db.CONNECTION_URI + db.CONNECTION_PORT)
}
