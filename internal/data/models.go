package data

import (
	"database/sql"

	"github.com/google/uuid"
)

type Models struct {
	User UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		User: UserModel{DB: db},
	}
}

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
}
