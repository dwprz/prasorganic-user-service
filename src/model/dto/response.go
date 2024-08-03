package dto

import "time"

type UserWithCredentials struct {
	UserId       string    `json:"user_id"`
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
	UserId       string    `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	PhotoProfile string    `json:"photo_profile"`
	Whatsapp     string    `json:"whatsapp"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
