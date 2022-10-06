package user

import "time"

type User struct {
	ID         int
	Name       string
	Occupation string
	Email      string `gorm:"index:idx_email,unique"`
	Password   string
	Avatar     string
	Role       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}