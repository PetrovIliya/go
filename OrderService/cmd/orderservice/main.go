package main

import (
	"OrderService/pkg/orderservice"
	"OrderService/pkg/repository"
	"OrderService/pkg/repositoryManager"
	"context"
	_ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)
import log "github.com/sirupsen/logrus"

func main() {
	config, err := ParseEnv()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	db, err := sqlx.Open("mysql", config.DataBaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sqlx.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(db)

	orderRepository := repository.OrderRepository{Db: *db}
	rm := repositoryManager.Create(orderRepository)
	server := orderservice.Server{RepositoryManager: rm}

	if err != nil {
		log.Error(err.Error())
		return
	}
	serverUrl := config.ServeRESTAddress
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting the server")
	killSignalChan := getKillSignalChan()
	httpServer := startHttpServer(serverUrl, &server)

	waitForKillSignal(killSignalChan)
	_ = httpServer.Shutdown(context.Background())
}

func startHttpServer(serverUrl string, server *orderservice.Server) *http.Server {
	router := orderservice.Router(server)
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
