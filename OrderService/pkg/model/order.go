package model

type Order struct {
	Id int `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
	CreatedAt string `json:"createdAt"`
	Cost float64 `json:"cost"`
}

type orderRepository interface {
	GetById(id int) (Order, error)
	GetAll() ([]Order, error)
	Add(order Order) (int, error)
	Delete(id int) error
	Update(order Order) error
	isExist(id int) (bool, error)
}