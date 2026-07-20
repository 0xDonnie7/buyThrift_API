package data

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) InsertUser(user *User) error {
	query := `INSERT INTO users (id, email, password_hash, role)
		VALUES ($1, $2)
	`
	args := []any{user.ID, user.Email, user.PasswordHash, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := u.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
