package model

type Order struct {
	Id        int        `json:"id"`
	MenuItems []MenuItem `json:"menuItems"`
	CreatedAt string     `json:"createdAt"`
	Cost      float64    `json:"cost"`
}
