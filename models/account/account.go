package account

type Account struct {
	ID       int64  `json:"id"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Address  string `json:"address" validate:"required"`
}
