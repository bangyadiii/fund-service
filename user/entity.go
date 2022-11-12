package user

import (
	"time"
)

type User struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	Occupation string
	Email      string `gorm:"index:idx_email,unique"`
	Password   string `json:"-"`
	Avatar     string
	Role       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
