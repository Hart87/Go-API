package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hart87/go-api/handlers"
	"github.com/hart87/go-api/models"
)

//Integration tests involving handlers and database

func TestCreateUser(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/users/new",
		strings.NewReader(`{"name":"Testy McTest","email":"test@gmail.com","password":"password","membership":"standard"}`))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.NewUserRoute)
	handler.ServeHTTP(rr, req)

	log.Print(rr.Body)
	var user models.User
	bodyBytes, _ := ioutil.ReadAll(rr.Body)
	err := json.Unmarshal(bodyBytes, &user)
	if err != nil {
		log.Print(err.Error())
	}

	os.Setenv("USER_ID", user.ID)
	log.Print(os.Getenv("USER_ID"))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestLoginForJWT(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/login",
		strings.NewReader(`{"email":"test@gmail.com","password":"password"}`))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.LoginRoute)
	handler.ServeHTTP(rr, req)

	log.Print(rr.Body)
	os.Setenv("TOKEN", rr.Body.String())
	log.Print(os.Getenv("TOKEN"))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestGetAllForResponseAndBody(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/users/all",
		nil)
	req.Header.Set("Token", os.Getenv("TOKEN"))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.GetAll)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	log.Print(rr.Body.String())

}

func TestGetById(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/users/"+os.Getenv("USER_ID"),
		nil)
	req.Header.Set("Token", os.Getenv("TOKEN"))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.UsersRoute)
	handler.ServeHTTP(rr, req)

	log.Print(rr.Body.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdateUser(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodPut,
		"/v1/users/"+os.Getenv("USER_ID"),
		strings.NewReader(`{"name":"Testy Testinson McTest IV","id":"`+os.Getenv("USER_ID")+`"}`))
	req.Header.Set("Token", os.Getenv("TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	log.Print(os.Getenv("TOKEN"))
	log.Print(os.Getenv("USER_ID"))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.UsersRoute)
	handler.ServeHTTP(rr, req)

	log.Print(rr.Body)
	var user models.User
	bodyBytes, _ := ioutil.ReadAll(rr.Body)
	err := json.Unmarshal(bodyBytes, &user)
	if err != nil {
		log.Print(err.Error())
	}
	log.Print("USERS NAME : " + user.Name)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDeleteById(t *testing.T) {

	req := httptest.NewRequest(
		http.MethodDelete,
		"/v1/users/"+os.Getenv("USER_ID"),
		nil)
	req.Header.Set("Token", os.Getenv("TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	log.Print(os.Getenv("TOKEN"))

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(handlers.UsersRoute)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
