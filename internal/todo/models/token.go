package models

type UserToken struct {
	ID           uint
	UserID       uint
	RefreshToken string
}
