package main

import (
	"log"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-send-raw-tx-tool/server/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		log.Fatal(err)
	}
}

func main() {

	r := gin.New()

	r.Use(cors.Default())

	gin.SetMode(gin.ReleaseMode)

	api := r.Group("/api/v1")

	api.GET("/info", handlers.GetInfo)

	api.POST("/send", handlers.Send)

	// TON testnet

	//api.POST("/sendGrams", handlers.SendGrams)
	//
	//api.POST("/signMessageHash", handlers.SigningMessageHash)

	if err := r.Run(":80"); err != nil {
		log.Fatal(err)
	}
}
