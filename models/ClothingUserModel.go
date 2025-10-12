package models

import (
	"time"
)

type ClothingUser struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Pin        int       `json:"pin"`
	UserStatus int       `json:"user_status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
