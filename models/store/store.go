package store

type Store struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userID"`
	NameStore   string `json:"nameStore" validate:"required"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdateAt    string `json:"update_at"`
}
