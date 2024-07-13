package models

type UserResponse struct {
	Id        int32  `db:"id"`
	Username  string `db:"username"`
	Email     string `db:"email"`
	Role      string `db:"role"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}
