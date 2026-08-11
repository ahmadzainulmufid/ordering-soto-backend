package models

import (
	"time"
)

type Category struct {
	ID           int64          `db:"id"`
	Name	     string         `db:"name"`
	Slug         string         `db:"slug"`
	IsActive     bool           `db:"is_active"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}