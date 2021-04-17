package orderservice

import (
	"OrderService/pkg/model"
	"OrderService/pkg/repositoryManager"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	RepositoryManager repositoryManager.RepositoryManager
}

func Router(srv *Server) http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/orders", srv.getOrders)
	s.HandleFunc("/order/{orderId}", srv.getOrder).Methods(http.MethodGet)
	s.HandleFunc("/order/{orderId}", srv.deleteOrder).Methods(http.MethodDelete)
	s.HandleFunc("/order/{orderId}", srv.updateOrder).Methods(http.MethodPut)
	s.HandleFunc("/order", srv.createOrder).Methods(http.MethodPost)

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

func (srv Server) getOrders(w http.ResponseWriter, _ *http.Request) {
	orderRepository := srv.RepositoryManager.GetOrderRepository()
	orders, err := orderRepository.GetAllOrders()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (srv Server) getOrder(w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	orderId, err := strconv.Atoi(requestParams["orderId"])
	if err != nil {
		http.Error(w, `Order id should contain only numbers [0-9]`, http.StatusInternalServerError)
		return
	}
	orderRepository := srv.RepositoryManager.GetOrderRepository()
	order, err := orderRepository.GetOrderById(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (srv Server) createOrder(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(r.Body)

	var order model.Order
	err = json.Unmarshal(b, &order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	menuItems := order.MenuItems
	if len(menuItems) == 0 {
		http.Error(w, "Wrong quantity of items", http.StatusBadRequest)
	}
	orderRepository := srv.RepositoryManager.GetOrderRepository()
	orderId, err := orderRepository.CreateOrder(order)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, `{ "orderId":` + strconv.Itoa(orderId) + ` }`); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Write Response Error")
	}
}

func (srv Server) deleteOrder (w http.ResponseWriter, r *http.Request) {
	requestParams := mux.Vars(r)
	orderId, err := strconv.Atoi(requestParams["orderId"])
	if err != nil {
		http.Error(w, `Order id should contain only numbers [0-9]`, http.StatusInternalServerError)
		return
	}
	orderRepository := srv.RepositoryManager.GetOrderRepository()
	isOrderExist, err := orderRepository.IsOrderExist(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isOrderExist {
		http.Error(w, `Order is not exist`, http.StatusBadRequest)
		return
	}
	err = orderRepository.DeleteOrderByOrderId(orderId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, `{ "success": "true" }`); err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Write Response Error")
	}
}

func (srv Server) updateOrder (w http.ResponseWriter, r *http.Request) {

}