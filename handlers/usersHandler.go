package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hart87/go-api/auth"
	"github.com/hart87/go-api/models"
)

//handler to base stuff off of for now....
func HandleRoute(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Name:       "John McGillicuddy",
		Email:      "JMC@gmail.com",
		Password:   "password",
		ID:         uuid.NewString(),
		Membership: "standard",
		CreatedAt:  1351700038}

	message, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	log.Print("test handler")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(message))
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
		log.Print("FORBIDDEN")
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
	time.Sleep(time.Second * 7) //Timeout Test
	log.Print("Post a User")
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("POST"))
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
