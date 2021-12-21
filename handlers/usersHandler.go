package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hart87/go-api/auth"
	"github.com/hart87/go-api/db"
	"github.com/hart87/go-api/models"

	"github.com/go-redis/redis"
	cache "github.com/hart87/go-api/redis"

	"github.com/golang-jwt/jwt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Redis
var rClient = redis.NewClient(&redis.Options{
	Addr:     cache.CONNECTION_URI + cache.CONNECTION_PORT,
	Password: cache.PASSWORD,
	DB:       cache.DB,
})

//Mongo

func GetAllUsers(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection, client, err := db.GetMongoDbCollection(db.DATABASE, db.COLLECTION_USERS)
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
	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	mResults, err := json.Marshal(results)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	client.Disconnect(ctx)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mResults))
}

func UsersRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getUserById(w, r)
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

func NewUserRoute(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		postUser(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
		return
	}
}

func getUserById(w http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 4 {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	part := parts[3]

	//check if Redis has the user before hitting Mongo
	rVal, cacheGetErr := rClient.Get(part).Result()
	if cacheGetErr == nil {
		log.Print("retrieved : " + rVal)
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rVal))
		return
	} else {
		log.Println("Not Presently Cache'd")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection, client, err := db.GetMongoDbCollection(db.DATABASE, db.COLLECTION_USERS)
	if err != nil {
		log.Panic(err)
	}

	result := models.User{}

	filter := bson.D{{"id", part}}
	val := collection.FindOne(ctx, filter).Decode(&result)
	if val != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader((http.StatusInternalServerError))
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(result)

	//Add to cache
	cacheSetError := rClient.Set(part, response, 0).Err()
	if cacheSetError != nil {
		log.Println("Not Cached : " + cacheSetError.Error())
	}

	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	client.Disconnect(ctx)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func postUser(w http.ResponseWriter, r *http.Request) {

	log.Println("POST USER")

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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user.CreatedAt = 1351700038
	user.ID = uuid.NewString()
	user.Password = auth.HashPassword(user.Password)
	log.Println(user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection, client, err := db.GetMongoDbCollection(db.DATABASE, db.COLLECTION_USERS)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println(res)

	jres, _ := json.Marshal(user)

	cacheSetError := rClient.Set(user.ID, jres, 0).Err()
	if cacheSetError != nil {
		log.Println("Not Cached : " + cacheSetError.Error())
	}

	client.Disconnect(ctx)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jres))
}

func editAUserById(w http.ResponseWriter, r *http.Request) {

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 4 {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	part := parts[3]

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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection, client, err := db.GetMongoDbCollection(db.DATABASE, db.COLLECTION_USERS)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"id", part}}
	update := bson.D{{"$set", user}}

	result, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if result.MatchedCount != 0 {
		log.Println("matched and replaced an existing document")
	}
	if result.UpsertedCount != 0 {
		log.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	res, _ := json.Marshal(user)

	cacheSetError := rClient.Set(part, res, 0).Err()
	if cacheSetError != nil {
		log.Println("Not Cached : " + cacheSetError.Error())
	}

	client.Disconnect(ctx)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func deleteAUserById(w http.ResponseWriter, r *http.Request) {

	type MyCustomClaims struct {
		ID   string `json:"id"`
		Role string `json:"role"`
		jwt.StandardClaims
	}

	//Obtain & parse token
	token, err := jwt.ParseWithClaims(r.Header["Token"][0], &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			w.WriteHeader(http.StatusForbidden)
			return nil, fmt.Errorf("something went wrong") //work on this line
		}
		return mySigningKey, nil
	})

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 4 {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	part := parts[3]

	//Obtain claims from token
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("claims not properly parsed from token"))
	}

	//Match conditions and possibly proceed
	log.Print(part)        //test purposes
	log.Print(claims.ID)   //test purposes
	log.Print(claims.Role) //test purposes

	if claims.Role != "admin" && part != claims.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Permission is not granted to delete this entity"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection, client, err := db.GetMongoDbCollection(db.DATABASE, db.COLLECTION_USERS)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	cacheDelError := rClient.Del(part).Err()
	if cacheDelError != nil {
		log.Println("Not Cached : " + cacheDelError.Error())
	}

	opts := options.Delete()
	filter := bson.D{{"id", part}}
	res, err := collection.DeleteOne(context.TODO(), filter, opts)
	if err != nil {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	log.Println(res)

	client.Disconnect(ctx)
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
}
