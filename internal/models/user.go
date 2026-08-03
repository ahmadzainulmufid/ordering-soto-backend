package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64          `db:"id"`
	FullName     string         `db:"full_name"`
	Email        string         `db:"email"`
	Phone        sql.NullString `db:"phone"`
	PasswordHash string         `db:"password_hash"`
	Role         string         `db:"role"`
	IsActive     bool           `db:"is_active"`
	LastLoginAt  sql.NullTime   `db:"last_login_at"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}