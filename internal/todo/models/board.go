package models

import "time"

type Board struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
