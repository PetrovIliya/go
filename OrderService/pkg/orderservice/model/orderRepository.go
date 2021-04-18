package model

type OrderRepository interface {
	GetById(id int) (Order, error)
	GetAll() ([]Order, error)
	Add(order Order) (int, error)
	Delete(id int) error
	Update(order Order) error
	IsExist(id int) (bool, error)
}
