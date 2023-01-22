package store

import "time"

type Store struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	NameStore   string    `json:"nameStore" validate:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"update_at"`
}
