package entities

import "time"

// domain Entity
type Users struct {
	ID          string    `gorm:"primarykey;autoIncrement" json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	AvatarURL   string    `json:"avatar_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
