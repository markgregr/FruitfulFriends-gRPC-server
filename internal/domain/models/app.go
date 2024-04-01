package models

type App struct {
	ID     int64  `gorm:"primaryKey"`
	Name   string `gorm:"unique"`
	Secret string `gorm:"not null"`
}
