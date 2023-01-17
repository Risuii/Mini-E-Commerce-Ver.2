package item

type Item struct {
	ID          int64  `json:"id"`
	StoreID     int64  `json:"userID"`
	Name        string `json:"nameStore" validate:"required"`
	Description string `json:"description"`
	Quantity    int64  `json:"quantity"`
}
