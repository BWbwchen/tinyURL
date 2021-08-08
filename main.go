package main

import (
	"github.com/gin-gonic/gin"
)

func init() {
	InitDatabase()
	InitRedis()
	InitZookeeper()
}

func main() {
	r := setupRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/add", AddURLHandler)
	r.GET("/:shortName", GetURLHandler)
	return r
}
