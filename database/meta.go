package database

import "time"

type Meta struct {
	UUID string `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
