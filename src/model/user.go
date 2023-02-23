package model

import (
	"time"
)

type User struct {
	ID           string `gorm:"primaryKey"`
	Name         string
	Occupation   string
	Email        string `gorm:"index:idx_email,unique"`
	Password     string `json:"-"`
	Avatar       string
	Role         string
	Transactions []Transaction
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
