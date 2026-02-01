package models

import "time"

type Poli struct {
	ID        uint      `gorm:"primaryKey"`
	NamaPoli  string    `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
