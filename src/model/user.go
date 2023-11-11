package model

import (
	"time"
)

type User struct {
	ID              string `gorm:"primaryKey"`
	Name            string
	Occupation      string
	Email           string `gorm:"index:idx_email,unique"`
	Password        string `json:"-"`
	Avatar          string
	Role            string `gorm:"default:user"`
	IsGoogleAccount bool   `json:"-"`
	Transactions    []Transaction
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
