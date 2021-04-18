package transport

import (
	"OrderService/pkg/orderservice/infrastructure"
	"OrderService/pkg/orderservice/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrders(t *testing.T) {
	const RequestAddress = "http://localhost:8000/api/v1/orders"
	req := httptest.NewRequest("GET", RequestAddress, nil)
	w := httptest.NewRecorder()

	srv := Server{orderRepository: infrastructure.CreateRepositoryMock(map[int]model.Order{})}
	srv.getOrders(w, req)
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have %d, wont %d", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	var orders []model.Order
	if err = json.Unmarshal(jsonString, &orders); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
	fmt.Println("Test passed with response: " + string(jsonString))
}

func TestGetOrder(t *testing.T) {
	const RequestAddress = "http://localhost:8000/api/v1/order/11"
	req := httptest.NewRequest("GET", RequestAddress, nil)
	w := httptest.NewRecorder()
	srv := Server{orderRepository: infrastructure.CreateRepositoryMock(map[int]model.Order{})}
	srv.getOrder(w, req)
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Have %d, wont %d", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	var order model.Order
	if err = json.Unmarshal(jsonString, &order); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
	fmt.Println("Test passed with response: " + string(jsonString))
}
