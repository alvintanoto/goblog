package model

import "time"

type Base struct {
	CreatedAt time.Time `gorm:"column:created_at"`
	CreatedBy int       `gorm:"column:created_by"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	UpdatedBy int       `gorm:"column:updated_by"`
}
