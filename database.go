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

type DBStruct struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type DB interface {
	GetURL(string) string
	GetShortName(string) string
	AddURLPair(string, string)
	URLExist(string) bool
	ShortNameExist(string) bool
}

var db DB

//GetMongoClient - Return mongodb connection to work with
func InitDatabase() {
	client := connectDB()
	_db := DBStruct{
		Client:     client,
		Collection: client.Database("url").Collection("url"),
	}
	db = _db
}

func connectDB() *mongo.Client {
	var mongoOnce sync.Once
	var client *mongo.Client
	var clientInstanceError error
	var CONNECTIONSTRING = "mongodb://" + os.Getenv("DB_URL")
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		// Set client options
		clientOptions := options.Client().ApplyURI(CONNECTIONSTRING)
		// Connect to MongoDB
		_client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = _client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		client = _client
	})
	if clientInstanceError != nil {
		panic("DB connection error")
	}
	return client
}

type data struct {
	LongURL   string `json:"longurl"`
	ShortName string `json:"shortname"`
}

func (dbs DBStruct) AddURLPair(LongURL string, shortName string) {
	toInsert := data{LongURL: LongURL, ShortName: shortName}
	_, err := dbs.Collection.InsertOne(context.TODO(), toInsert)
	if err != nil {
		fmt.Println("Insert Error")
	}
}

func (dbs DBStruct) GetURL(shortName string) string {
	result, _ := dbs.find("shortname", shortName)
	return result.LongURL
}

func (dbs DBStruct) URLExist(LongURL string) bool {
	_, err := dbs.find("longurl", LongURL)
	return err == nil
}

func (dbs DBStruct) ShortNameExist(shortName string) bool {
	_, err := dbs.find("shortname", shortName)
	return err == nil
}

func (dbs DBStruct) GetShortName(LongURL string) string {
	result, _ := dbs.find("longurl", LongURL)
	return result.ShortName
}

func (dbs DBStruct) find(columnName, value string) (data, error) {
	var result data
	err := dbs.Collection.FindOne(context.TODO(), bson.M{columnName: value}).Decode(&result)
	return result, err
}
