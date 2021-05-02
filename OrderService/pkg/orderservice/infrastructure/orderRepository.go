package infrastructure

import (
	model "OrderService/pkg/orderservice/model"
	"database/sql"
)

func CreateRepository(db *sql.DB) model.OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

type OrderRepository struct {
	db *sql.DB
}

func (repository OrderRepository) GetById(id int) (model.Order, error) {
	var order model.Order
	var menuItems []model.MenuItem
	err := repository.db.QueryRow("SELECT order_id AS id, cost, created_at AS createdAt FROM `order` WHERE order_id = ?", id).Scan(&order.Id, &order.Cost, &order.CreatedAt)
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
	rows, err := repository.db.Query("SELECT order_id AS id, cost, created_at AS createdAt FROM `order`")
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

func (repository OrderRepository) Add(order model.Order) (int, error) {
	menuItems := order.MenuItems
	orderCost := order.Cost

	res, err := repository.db.Exec("INSERT INTO `order` (cost) VALUES(?)", orderCost)
	if err != nil {
		return 0, err
	}
	orderId, _ := res.LastInsertId()

	err = repository.createMenuItems(menuItems, int(orderId))
	return int(orderId), nil
}

func (repository OrderRepository) IsExist(orderId int) (bool, error) {
	count := 0
	err := repository.db.QueryRow("SELECT COUNT(*) AS count FROM `order` WHERE order_id = ?", orderId).Scan(&count)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (repository OrderRepository) Delete(orderId int) error {
	_, err := repository.db.Exec("DELETE FROM `order` WHERE order_id = ?", orderId)
	if err != nil {
		return err
	}
	err = repository.deleteOrderMenuItems(orderId)
	return nil
}

func (repository OrderRepository) Update(order model.Order) error {
	menuItems := order.MenuItems

	_, err := repository.db.Exec("UPDATE `order` SET cost = ?, updated_at = NOW() WHERE order_id = ?", order.Cost, order.Id)
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

func (repository OrderRepository) getMenuItemsByOrderId(orderId int) ([]model.MenuItem, error) {
	var menuItems []model.MenuItem
	rows, err := repository.db.Query("SELECT menu_item_id AS id, name, quantity FROM menu_item WHERE order_id = ?", orderId)
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
	_, err := repository.db.Exec("DELETE FROM menu_item WHERE order_id = ?", orderId)
	if err != nil {
		return err
	}

	return nil
}

func (repository OrderRepository) createMenuItems(menuItems []model.MenuItem, orderId int) error {
	for i := 0; i < len(menuItems); i++ {
		menuItem := menuItems[i]
		_, err := repository.db.Exec("INSERT INTO menu_item (name, quantity, order_id) VALUES (?, ?, ?)", menuItem.Name, menuItem.Quantity, orderId)
		if err != nil {
			return err
		}
	}
	return nil
}
