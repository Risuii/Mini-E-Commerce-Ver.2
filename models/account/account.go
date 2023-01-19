package account

import "time"

type Account struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	Email     string    `json:"email" validate:"email"`
	Address   string    `json:"address" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
}
