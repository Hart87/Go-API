package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hart87/go-api/auth"
	"github.com/hart87/go-api/db"
	"github.com/hart87/go-api/models"

	"go.mongodb.org/mongo-driver/bson"
)

func TestRoute(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Name:       "John McGillicuddy",
		Email:      "JMC@gmail.com",
		Password:   auth.HashPassword("password"),
		ID:         uuid.NewString(),
		Membership: "standard",
		CreatedAt:  1351700038}

	message, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	//hash test
	log.Println(user.Password)
	log.Println(auth.CheckPasswordHash(user.Password, "password"))

	log.Print("test handler")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if !auth.IsAuthenticated(r) {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - FORBIDDEN"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection, err := db.GetMongoDbCollection(db.DATABASE, db.COLLECTION_USERS)
	if err != nil {
		log.Panic(err)
	}

	filter := bson.D{}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}

	defer cursor.Close(ctx)
	results := []models.User{}

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	mResults, err := json.Marshal(results)

	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mResults))
}

func UsersRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUserById(w, r)
		return
	case "POST":
		postUser(w, r)
		return
	case "PUT":
		editAUserById(w, r)
		return
	case "DELETE":
		deleteAUserById(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	//Authentication here
	if !auth.IsAuthenticated(r) {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 - FORBIDDEN"))
		return
	}

	log.Print("Get a User by Id")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GET"))
}

func postUser(w http.ResponseWriter, r *http.Request) {

	//Dummy user here for a test
	user := models.User{
		Name:       "Steve Herman",
		Email:      "StevieH@yahooo.com",
		Password:   auth.HashPassword("password"),
		ID:         uuid.NewString(),
		Membership: "admin",
		CreatedAt:  13517005000}

	//db.CreateUser(user)
	result, _ := json.Marshal(user)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func editAUserById(w http.ResponseWriter, r *http.Request) {
	//Authentication here
	log.Print("Edit a User")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PUT"))
}

func deleteAUserById(w http.ResponseWriter, r *http.Request) {
	//Authentication here
	log.Print("Delete a User")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("DELETE"))
}
