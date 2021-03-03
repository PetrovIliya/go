package main

import (
	"OrderService/pkg/orderservice"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import log "github.com/sirupsen/logrus"

func main() {
	const LogFile = "var/log/myLog.log"
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}

	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	killSignalChan := getKillSignalChan()
	srv := startServer(serverUrl)

	waitForKillSignal(killSignalChan)
	_ = srv.Shutdown(context.Background())
}

func startServer(serverUrl string) *http.Server {
	router := orderservice.Router()
	srv := &http.Server{Addr: serverUrl, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	return srv
}

func getKillSignalChan() chan os.Signal {
	osKillSignalChan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalChan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalChan
}

func waitForKillSignal(getKillSignalChan <-chan os.Signal) {
	killSignal := <-getKillSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM...")
	}
}
