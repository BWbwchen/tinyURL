package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type POSTData struct {
	LongURL string `json:"URL"`
}

type POSTresponse struct {
	Status    string `json:"status"`
	ShortName string `json:"short"`
}

func AddURLHandler(c *gin.Context) {
	var requestBody POSTData
	c.ShouldBindJSON(&requestBody)

	var shortName string
	if db.URLExist(requestBody.LongURL) {
		shortName = db.GetShortName(requestBody.LongURL)
	} else {
		shortName = GenerateShortName(requestBody.LongURL)
		storeIntoDBRedis(requestBody.LongURL, shortName)
	}

	c.JSON(http.StatusOK, POSTresponse{
		Status:    "Exist",
		ShortName: shortName,
	})
}

func storeIntoDBRedis(LongURL, shortName string) {
	var wg sync.WaitGroup
	wg.Add(2)
	// add in database
	go func(LongURL string, shortName string) {
		defer wg.Done()
		db.AddURLPair(LongURL, shortName)
	}(LongURL, shortName)
	// add in redis
	go func(LongURL string, shortName string) {
		defer wg.Done()
		RedisAdd(LongURL, shortName)
	}(LongURL, shortName)

	wg.Wait()
}
