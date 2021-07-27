package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	InitDatabase()
	InitRedis()
	InitZookeeper()
}

type POSTData struct {
	LongURL string `json:"URL"`
}

type POSTresponse struct {
	Status    string `json:"status"`
	ShortName string `json:"short"`
}
type GETresponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

func main() {
	r := gin.Default()
	r.POST("/add", func(c *gin.Context) {
		var body POSTData
		c.ShouldBindJSON(&body)
		// validate it
		// generate the corresponding short name
		shortName, isExist := DatabaseURLExist(body.LongURL)
		if isExist {
			c.JSON(http.StatusOK, POSTresponse{
				Status:    "Exist",
				ShortName: shortName,
			})
		} else {
			shortName := GetShortName(body.LongURL)
			// add in database
			DatabaseAdd(body.LongURL, shortName)
			// add in redis
			RedisAdd(body.LongURL, shortName)

			c.JSON(http.StatusOK, POSTresponse{
				Status:    "OK",
				ShortName: shortName,
			})
		}
	})
	r.GET("/:shortName", func(c *gin.Context) {
		shortName := c.Param("shortName")
		// find in redis first
		LongURL, err := RedisGet(shortName)
		//LongURL, err := DatabaseGet(shortName)
		if err != nil {
			// if not found, find in database
			DatabaseLongURL, err_ := DatabaseGet(shortName)
			if err_ != nil {
				c.JSON(http.StatusNotFound, GETresponse{
					Status: "Not OK",
					URL:    "",
				})
			} else {
				RedisAdd(DatabaseLongURL, shortName)
				c.JSON(http.StatusOK, GETresponse{
					Status: "(Database)OK",
					URL:    DatabaseLongURL,
				})
			}
		} else {
			c.JSON(http.StatusOK, GETresponse{
				Status: "(Redis)OK",
				URL:    LongURL,
			})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
