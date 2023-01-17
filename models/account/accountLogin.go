package account

type AccountLogin struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}
