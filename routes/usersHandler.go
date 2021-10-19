package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hart87/go-api/models"
)

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
