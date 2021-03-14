package repository

import (
	"database/sql"
	_ "database/sql"
	"github.com/jmoiron/sqlx"
)

type Order struct {
	Id int `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
	CreatedAt string `json:"createdAt"`
	Cost float64 `json:"cost"`
}

type Orders struct {
	Orders []Order
}

type MenuItem struct {
	Id string `json:"id"`
	Quantity int `json:"quantity"`
	Name string `json:"name"`
}

type OrderRepository struct {
	Db sqlx.DB
}

func (repository OrderRepository) GetOrderById(orderId int) (Order, error) {
	var order Order
	var menuItems []MenuItem
	err := repository.Db.QueryRow("SELECT order_id AS id, cost, created_at AS createdAt FROM `order` WHERE order_id = ?", orderId).Scan(&order.Id, &order.Cost, &order.CreatedAt)
	if err != nil {
		return order, err
	}
	menuItems, err = repository.getMenuItemsByOrderId(orderId)
	if err != nil {
		return order, err
	}
	order.MenuItems = menuItems

	return order, nil
}

func (repository OrderRepository) GetAllOrders() (Orders, error) {
	var orders Orders
	rows, err := repository.Db.Query("SELECT order_id AS id, cost, created_at AS createdAt FROM `order`")
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if err != nil {
		return orders, err
	}
	for rows.Next() {
		var menuItems []MenuItem
		var order Order
		err := rows.Scan(&order.Id, &order.Cost, &order.CreatedAt)
		if err != nil {
			return orders, err
		}
		menuItems, err = repository.getMenuItemsByOrderId(order.Id)
		if err != nil {
			return orders, err
		}
		order.MenuItems = menuItems
		orders.Orders = append(orders.Orders, order)
	}
	return orders, nil
}

func (repository OrderRepository) CreateOrder(order Order) (int, error) {
	menuItems := order.MenuItems
	orderCost := order.Cost

	res, err := repository.Db.Exec("INSERT INTO `order` (cost) VALUES(?)", orderCost)
	if err != nil {
		return 0, err
	}
	orderId, _ := res.LastInsertId()

	err = repository.createMenuItems(menuItems, int(orderId))
	return int(orderId), nil
}

func (repository OrderRepository) IsOrderExist(orderId int) (bool, error) {
	count := 0
	err := repository.Db.QueryRow("SELECT COUNT(*) AS count FROM `order` WHERE order_id = ?", orderId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (repository OrderRepository) DeleteOrderByOrderId(orderId int) error {
	_, err := repository.Db.Exec("DELETE FROM `order` WHERE order_id = ?", orderId)
	if err != nil{
		return err
	}
	err = repository.deleteOrderMenuItems(orderId)
	return nil
}

func (repository OrderRepository) UpdateOrder(order Order) error {
	menuItems := order.MenuItems

	_, err := repository.Db.Exec("UPDATE `order` SET cost = ?, update_at = NOW() WHERE order_id = ?", order.Cost, order.Id)
	if err != nil {
		return err
	}
	err = repository.deleteOrderMenuItems(order.Id)
	if err != nil {
		return err
	}
	err = repository.createMenuItems(menuItems, order.Id)
	return nil
}

func (repository OrderRepository) getMenuItemsByOrderId(orderId int) ([]MenuItem, error)  {
	var menuItems []MenuItem
	rows, err := repository.Db.Query("SELECT menu_item_id AS id, name, quantity FROM menu_item WHERE order_id = ?", orderId)
	if err != nil {
		return menuItems, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		var menuItem MenuItem
		err := rows.Scan(&menuItem.Id, &menuItem.Name, &menuItem.Quantity)
		if err != nil {
			return menuItems, err
		}
		menuItems = append(menuItems, menuItem)
	}

	return menuItems, nil
}

func (repository OrderRepository) deleteOrderMenuItems(orderId int) error {
	_, err := repository.Db.Exec("DELETE FROM menu_item WHERE order_id = ?", orderId)
	if err != nil{
		return err
	}

	return nil
}

func (repository OrderRepository) createMenuItems(menuItems []MenuItem, orderId int) error {
	for i := 0; i < len(menuItems); i++ {
		menuItem := menuItems[i]
		_, err := repository.Db.Exec("INSERT INTO menu_item (name, quantity, order_id) VALUES (?, ?, ?)", menuItem.Name, menuItem.Quantity, orderId)
		if err != nil {
			return err
		}
	}
	return nil
}



