package dto

import "time"

type UserCreate struct {
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type UserWithCredentials struct {
	UserID       uint      `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	PhotoProfile string    `json:"photo_profile"`
	Whatsapp     string    `json:"whatsapp"`
	Password     string    `json:"password"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserSanitized struct {
	UserID       uint      `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	PhotoProfile string    `json:"photo_profile"`
	Whatsapp     string    `json:"whatsapp"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
