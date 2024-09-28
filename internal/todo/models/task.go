package models

import "time"

type Task struct {
	ID          uint      `gorm:"primeryKey"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	Board_id    uint      `gorm:"column:board_id"`
	Status_id   uint      `gorm:"column:status_id"`
	User_id     uint      `gorm:"column:user_id"`
	Created_at  time.Time `gorm:"column:created_at"`
	Updated_at  time.Time `gorm:"column:updated_at"`
}
