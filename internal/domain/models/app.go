package models

type App struct {
	ID     int64  `gorm:"primaryKey" index:"idx_id" json:"id"`
	Name   string `gorm:"unique" index:"idx_name" json:"name"`
	Secret string `gorm:"not null" json:"secret"`
}
