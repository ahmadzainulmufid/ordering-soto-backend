package models

import (
	"database/sql"
	"time"
)

type Product struct {
	ID		  int64     `db:"id"`
	CategoryID int64     `db:"category_id"`
	Name       string    `db:"name"`
	Slug       string    `db:"slug"`
	Description sql.NullString `db:"description"`
	Price      float64   `db:"price"`
	ImageURL   sql.NullString `db:"image_url"`
	Stock	  int       `db:"stock"`
	IsAvailable  bool      `db:"is_available"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}