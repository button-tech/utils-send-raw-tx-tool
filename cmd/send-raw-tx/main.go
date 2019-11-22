package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/button-tech/utils-send-raw-tx-tool/api"
)

func main() {
	server := api.NewServer()

	signalEx := make(chan os.Signal, 1)
	defer close(signalEx)

	signal.Notify(signalEx,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		log.Println("start http server")
		if err := server.Core.ListenAndServe(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	defer func() {
		if err := server.Core.Shutdown(); err != nil {
			log.Fatal(err)
		}
	}()

	stop := <-signalEx
	log.Println("Received", stop)
	log.Println("Waiting for all jobs to stop")
}
