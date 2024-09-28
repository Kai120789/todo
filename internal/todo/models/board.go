package models

import "time"

type Board struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"column:name"`
	User_id    uint      `gorm:"column:user_id"`
	Created_at time.Time `gorm:"column:created_at"`
}
