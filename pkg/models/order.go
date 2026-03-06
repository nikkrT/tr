package models

import "time"

type OrderModel struct {
	Id        int    `json:"id"`
	ProductId int    `json:"product_id"`
	Status    string `json:"status"`
	CreatedAt time.Time
	DeletedAt time.Time
}
