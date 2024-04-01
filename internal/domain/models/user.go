package models

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	PassHash []byte `gorm:"not null"`
	Role     int    `gorm:"not null"`
	Status   int    `gorm:"not null"`
}
