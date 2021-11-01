package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hart87/go-api/handlers"
)

//Integration tests involving handler and database

func TestGetAllForResponseAndBody(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/all",
		nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.GetAllUsers)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	log.Print(rr.Body.String())

}

func TestGetByIdForResponseAndBody(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/261f0214-1097-4c2b-ae3a-a86b15f25e6b",
		nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.UsersRoute)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"name":"Barney Stinson","email":"BStin@hotmail.com","password":"$2a$10$A1tAnFJ2voUruJYpCqW0uekV3fqkfd7I8MfAbYU82trEAGEBkKVUq","id":"261f0214-1097-4c2b-ae3a-a86b15f25e6b","membership":"basic","createdAt":1351700038}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
