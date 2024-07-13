package models

type OrderItem struct {
	ProductId int32   `db:"product_id"`
	Quantity  int32   `db:"quantity"`
	Price     float32 `db:"price"`
}

type OrderResponse struct {
	Id          int32       `db:"id"`
	UserId      int32       `db:"user_id"`
	Items       []OrderItem `db:"items"`
	TotalAmount float32     `db:"total_amount"`
	Status      string      `db:"status"`
	CreatedAt   string      `db:"created_at"`
	UpdatedAt   string      `db:"updated_at"`
}
