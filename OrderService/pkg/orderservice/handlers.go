package orderservice

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type Orders struct {
	Orders []Order
}

type MenuItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	Id                 string     `json:"id"`
	MenuItems          []MenuItem `json:"menuItems"`
	OrderedAtTimesTime int        `json:"orderedAtTimesTime"`
	Cost               float64    `json:"cost"`
}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-word", helloWord)
	s.HandleFunc("/orders", getOrders)
	s.HandleFunc("/order/{orderId}", getOrder)

	return logMiddleWare(r)
}

func logMiddleWare(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		duration := time.Since(startTime)
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
			"duration":   duration.Seconds(),
		}).Info("got a new request")

	})
}

func helloWord(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "Hello word!")
}

func getOrders(w http.ResponseWriter, _ *http.Request) {
	menuItems := []MenuItem{{Id: "someMenuId", Quantity: 0}}
	order := Order{Id: "someOrderId", MenuItems: menuItems}
	orders := Orders{[]Order{order}}
	b, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, string(b)); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Write Response Error")
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	orderId := requestParams["orderId"]
	menuItems := []MenuItem{{Id: "someMenuId", Quantity: 0}}
	order := Order{Id: orderId, MenuItems: menuItems, Cost: 1, OrderedAtTimesTime: 22}
	b, err := json.Marshal(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, string(b)); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Write Response Error")
	}
}
