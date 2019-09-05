package main

import (
	"log"

	"github.com/button-tech/utils-send-raw-tx-tool/server/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	r.Use(cors.Default())

	gin.SetMode(gin.ReleaseMode)

	api := r.Group("/api/v1")

	api.GET("/info", handlers.GetInfo)

	api.POST("/send", handlers.Send)

	if err := r.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
