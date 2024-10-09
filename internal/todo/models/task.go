package models

import "time"

type Task struct {
	ID          uint
	Title       string
	Description string
	BoardId     uint
	StatusId    uint
	UserId      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
