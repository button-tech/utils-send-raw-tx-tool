package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/button-tech/logger"
	"github.com/button-tech/utils-send-raw-tx-tool/api"
)

const port = ":8080"

func main() {
	if err := logger.InitLogger(os.Getenv("DSN")); err != nil {
		logger.Fatal(err)
	}

	server, err := api.NewServer()
	if err != nil {
		logger.Fatal(err)
	}

	signalEx := make(chan os.Signal, 1)
	defer close(signalEx)

	signal.Notify(signalEx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		log.Printf("start http server on port:%s", port)
		if err := server.Core.ListenAndServe(port); err != nil {
			logger.Fatal(err)
		}
	}()

	defer func() {
		if err := server.Core.Shutdown(); err != nil {
			logger.Fatal(err)
		}
	}()

	stop := <-signalEx
	log.Println("Received", stop)
	log.Println("Waiting for all jobs to stop")
}
