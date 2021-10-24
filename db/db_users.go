package db

// import (
// 	"context"
// 	"log"

// 	"github.com/hart87/go-api/models"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

// func CreateUser(user models.User) {
// 	collection, err := getMongoDbCollection(DATABASE, COLLECTION_USERS)
// 	if err != nil {
// 		log.Panic(err)
// 	}
// 	collection.InsertOne(context.TODO(), user)
// }

// func FindAllUsers() []models.User{} {
// 	collection, err := getMongoDbCollection(DATABASE, COLLECTION_USERS)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	filter := bson.D{}
// 	cursor, err := collection.Find(context.TODO(), filter)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return users
// }
