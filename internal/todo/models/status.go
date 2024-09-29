package models

type Status struct {
	ID   uint   `gorm:"primaryKey"`
	Type string `gotm:"column:type"`
}
