package models

import "time"

type Poli struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    NamaPoli  string    `gorm:"unique;not null" json:"nama_poli"` 
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}