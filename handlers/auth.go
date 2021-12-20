package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/hart87/go-api/auth"
	"github.com/hart87/go-api/db"
	"github.com/hart87/go-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

var mySigningKey = []byte("pizza") //TEMPORARILY HERE

func Login(w http.ResponseWriter, r *http.Request) {

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var user models.User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection, client, err := db.GetMongoDbCollection(db.DATABASE, db.COLLECTION_USERS)
	if err != nil {
		log.Panic(err)
	}

	result := models.User{}

	//find the user from the request body
	filter := bson.D{{"email", user.Email}}
	val := collection.FindOne(ctx, filter).Decode(&result)
	if val != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader((http.StatusBadRequest))
		w.Write([]byte(err.Error()))
		return
	}

	//run a hash check on the passwords
	isOk := auth.CheckPasswordHash(result.Password, user.Password)
	if !isOk {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader((http.StatusUnauthorized))
		return
	}

	client.Disconnect(ctx)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader((http.StatusOK))
	w.Write([]byte(GenerateToken()))
	return
}

func GenerateToken() string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Print("something went wrong: %s", err.Error())
	}

	return tokenString
}

func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					w.WriteHeader(http.StatusForbidden)
					return nil, fmt.Errorf("something went wrong") //work on this line
				}
				return mySigningKey, nil
			})

			if err != nil {
				log.Print(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			log.Print(w, "Not Authorized")
			w.WriteHeader(http.StatusForbidden)
		}

	})

}
