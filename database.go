package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/* Used to create a singleton object of MongoDB client.
Initialized and exposed through  GetMongoClient().*/
var clientInstance *mongo.Client

//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

var collection *mongo.Collection

//I have used below constants just to hold required database config's.
var (
	CONNECTIONSTRING = "mongodb://" + os.Getenv("DB_URL")
)

type data struct {
	LongURL   string `json:"longurl"`
	ShortName string `json:"shortname"`
}

//GetMongoClient - Return mongodb connection to work with
func InitDatabase() {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	if clientInstanceError != nil {
		fmt.Println("Connection error")
	}
	collection = clientInstance.Database("url").Collection("url")
	//return clientInstance, clientInstanceError
}

func DatabaseAdd(LongURL string, shortName string) {
	toInsert := data{LongURL: LongURL, ShortName: shortName}
	res, err := collection.InsertOne(context.TODO(), toInsert)
	if err != nil {
		fmt.Println("Insert Error")
	}
	id := res.InsertedID
	fmt.Println("Insert " + fmt.Sprintf("%v", id) + " record !")
}

func DatabaseGet(shortName string) (string, error) {
	var result data

	err := collection.FindOne(context.TODO(), bson.M{"shortname": shortName}).Decode(&result)
	return result.LongURL, err
}

func DatabaseURLExist(LongURL string) (string, bool) {
	var result data

	err := collection.FindOne(context.TODO(), bson.M{"longurl": LongURL}).Decode(&result)
	return result.ShortName, (err == nil)
}
