package user

import (
	"backend-crowdfunding/transaction"
	"time"
)

type User struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Occupation   string
	Email        string `gorm:"index:idx_email,unique"`
	Password     string `json:"-"`
	Avatar       string
	Role         string
	Transactions []transaction.Transaction
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
