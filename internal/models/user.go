package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Email     string    `gorm:"uniqueIndex" json:"email"`
	Password  string    `json:"-"` // Don't return password in JSON
	UnitKerja string    `json:"unit_kerja"`
	TimKerja  string    `json:"tim_kerja"`
	Role      string    `json:"role"` // role admin, ketua tim kerja, ketua unitkerja, user
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
