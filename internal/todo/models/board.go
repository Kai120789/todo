package models

import "time"

type Board struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
