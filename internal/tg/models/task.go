package models

import "time"

type Task struct {
	ID          uint
	Title       string `json:"title"`
	Description string `json:"description"`
	BoardId     uint
	StatusId    uint `json:"status_id"`
	UserId      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
