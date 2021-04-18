package infrastructure

import (
	"OrderService/pkg/orderservice/model"
	"errors"
)

func CreateRepositoryMock(orders map[int]model.Order) model.OrderRepository {
	return &OrderRepositoryMock{
		orders: orders,
	}
}

type OrderRepositoryMock struct {
	orders map[int]model.Order
}

func (o OrderRepositoryMock) GetById(id int) (model.Order, error) {
	order, ok := o.orders[id]
	if !ok {
		return model.Order{}, errors.New("order is not exist")
	}
	return order, nil
}

func (o OrderRepositoryMock) GetAll() ([]model.Order, error) {
	return o.getOrdersAsArray(), nil
}

func (o OrderRepositoryMock) Add(order model.Order) (int, error) {
	o.orders[order.Id] = order
	return order.Id, nil
}

func (o OrderRepositoryMock) Delete(id int) error {
	delete(o.orders, id)
	return nil
}

func (o OrderRepositoryMock) Update(order model.Order) error {
	o.orders[order.Id] = order
	return nil
}

func (o OrderRepositoryMock) IsExist(id int) (bool, error) {
	_, isExist := o.orders[id]
	return isExist, nil
}

func (o OrderRepositoryMock) getOrdersAsArray() []model.Order {
	var orders []model.Order
	for _, element := range o.orders {
		orders = append(orders, element)
	}

	return orders
}
