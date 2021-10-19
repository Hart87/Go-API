package routes

import (
	"net/http"
)

func HandleRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
