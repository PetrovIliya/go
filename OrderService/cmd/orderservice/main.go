package main

import (
	"OrderService/pkg/orderservice/infrastructure"
	"OrderService/pkg/orderservice/transport"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import log "github.com/sirupsen/logrus"

func main() {
	config, err := ParseEnv()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	println(config.DatabaseUrl)

	db, err := sql.Open("mysql", config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(db)

	orderRepository := infrastructure.CreateRepository(db)
	server := transport.CreateServer(orderRepository)

	if err != nil {
		log.Error(err.Error())
		return
	}
	serverUrl := config.ServeRESTAddress
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	killSignalChan := getKillSignalChan()
	httpServer := startHttpServer(serverUrl, server)

	waitForKillSignal(killSignalChan)
	_ = httpServer.Shutdown(context.Background())
}

func startHttpServer(serverUrl string, server *transport.Server) *http.Server {
	router := transport.Router(server)
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
