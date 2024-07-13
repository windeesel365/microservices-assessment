package models

type GetProductResponse struct {
	Id          int32   `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Price       float32 `db:"price"`
	CategoryId  int32   `db:"category_id"`
	Stock       int32   `db:"stock"`
	CreatedAt   string  `db:"created_at"`
	UpdatedAt   string  `db:"updated_at"`
}
