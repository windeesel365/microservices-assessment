package models

type PaymentResponse struct {
	Id            int32   `db:"id"`
	OrderId       int32   `db:"order_id"`
	Amount        float32 `db:"amount"`
	PaymentMethod string  `db:"payment_method"`
	Status        string  `db:"status"`
	TransactionId string  `db:"transaction_id"`
	CreatedAt     string  `db:"created_at"`
}
