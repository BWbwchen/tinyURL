package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GETresponse struct {
	Status string `json:"status"`
	URL    string `json:"url"`
}

func GetURLHandler(c *gin.Context) {
	shortNameQuery := c.Param("shortName")

	LongURL, err := RedisGet(shortNameQuery)
	if err != nil {
		DatabaseLongURL := db.GetURL(shortNameQuery)
		if DatabaseLongURL == "" {
			c.JSON(http.StatusNotFound, GETresponse{
				Status: "Not OK",
				URL:    "",
			})
		} else {
			RedisAdd(DatabaseLongURL, shortNameQuery)
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
}
