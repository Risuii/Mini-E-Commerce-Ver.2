package item

import "time"

type Item struct {
	ID          int64     `json:"id"`
	StoreID     int64     `json:"userID"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	Quantity    int64     `json:"quantity" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
}
