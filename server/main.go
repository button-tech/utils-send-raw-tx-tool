package main

import (
	"github.com/button-tech/utils-send-raw-tx-tool/server/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	r := gin.New()

	gin.SetMode(gin.ReleaseMode)

	api := r.Group("/api/v1")

	api.GET("/info", handlers.GetInfo)

	api.POST("/send/:currency", handlers.Send)

	if err := r.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
