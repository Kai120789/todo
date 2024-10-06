package models

import "time"

type User struct {
	ID           uint
	Username     string
	TgName       string
	PasswordHash []byte
	CreatedAt    time.Time
}
