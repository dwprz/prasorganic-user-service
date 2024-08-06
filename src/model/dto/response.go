package dto

import (
	"time"
)

type SanitizedUserRes struct {
	UserId         string    `json:"user_id"`
	Email          string    `json:"email"`
	FullName       string    `json:"full_name"`
	Role           string    `json:"role"`
	PhotoProfileId string    `json:"photo_profile_id"`
	PhotoProfile   string    `json:"photo_profile"`
	Whatsapp       string    `json:"whatsapp"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpdateEmailRes struct {
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	NewEmail string `json:"new_email" validate:"required,email,min=5,max=100"`
}

type VerifyUpdateEmailRes struct {
	Data        *SanitizedUserRes
	AccessToken string
}
