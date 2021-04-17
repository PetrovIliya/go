package repository

import (
	"OrderService/pkg/model"
	"database/sql"
	_ "database/sql"
	"github.com/jmoiron/sqlx"
)

type OrderRepository struct {
	Db sqlx.DB
}

func (repository OrderRepository) GetById(id int) (model.Order, error) {
	var order model.Order
	var menuItems []model.MenuItem
	err := repository.Db.QueryRow("SELECT order_id AS id, cost, created_at AS createdAt FROM `order` WHERE order_id = ?", id).Scan(&order.Id, &order.Cost, &order.CreatedAt)
	if err != nil {
		return order, err
	}
	menuItems, err = repository.getMenuItemsByOrderId(id)
	if err != nil {
		return order, err
	}
	order.MenuItems = menuItems

	return order, nil
}

func (repository OrderRepository) GetAll() ([]model.Order, error) {
	var orders []model.Order
	rows, err := repository.Db.Query("SELECT order_id AS id, cost, created_at AS createdAt FROM `order`")
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	if err != nil {
		return orders, err
	}
	for rows.Next() {
		var menuItems []model.MenuItem
		var order model.Order
		err := rows.Scan(&order.Id, &order.Cost, &order.CreatedAt)
		if err != nil {
			return orders, err
		}
		menuItems, err = repository.getMenuItemsByOrderId(order.Id)
		if err != nil {
			return orders, err
		}
		order.MenuItems = menuItems
		orders = append(orders, order)
	}
	return orders, nil
}

func (repository OrderRepository) Insert(order model.Order) (int, error) {
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

func (repository OrderRepository) IsExist(orderId int) (bool, error) {
	count := 0
	err := repository.Db.QueryRow("SELECT COUNT(*) AS count FROM `order` WHERE order_id = ?", orderId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (repository OrderRepository) DeleteById(orderId int) error {
	_, err := repository.Db.Exec("DELETE FROM `order` WHERE order_id = ?", orderId)
	if err != nil{
		return err
	}
	err = repository.deleteOrderMenuItems(orderId)
	return nil
}

func (repository OrderRepository) UpdateById(id int, order model.Order) error {
	menuItems := order.MenuItems

	_, err := repository.Db.Exec("UPDATE `order` SET cost = ?, update_at = NOW() WHERE order_id = ?", order.Cost, id)
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

func (repository OrderRepository) getMenuItemsByOrderId(orderId int) ([]model.MenuItem, error)  {
	var menuItems []model.MenuItem
	rows, err := repository.Db.Query("SELECT menu_item_id AS id, name, quantity FROM menu_item WHERE order_id = ?", orderId)
	if err != nil {
		return menuItems, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	for rows.Next() {
		var menuItem model.MenuItem
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

func (repository OrderRepository) createMenuItems(menuItems []model.MenuItem, orderId int) error {
	for i := 0; i < len(menuItems); i++ {
		menuItem := menuItems[i]
		_, err := repository.Db.Exec("INSERT INTO menu_item (name, quantity, order_id) VALUES (?, ?, ?)", menuItem.Name, menuItem.Quantity, orderId)
		if err != nil {
			return err
		}
	}
	return nil
}



