package entity

import "time"

type User struct {
	UserId         string    `json:"user_id" gorm:"column:user_id;primaryKey"`
	Email          string    `json:"email" gorm:"column:email"`
	FullName       string    `json:"full_name" gorm:"column:full_name"`
	Role           string    `json:"role" gorm:"column:role;default:'USER'"`
	PhotoProfileId string    `json:"photo_profile_id" column:"photo_profile_id;default:null"`
	PhotoProfile   string    `json:"photo_profile" column:"photo_profile;default:null"`
	Whatsapp       string    `json:"whatsapp" gorm:"column:whatsapp;default:null"`
	Password       string    `json:"password" gorm:"column:password"`
	RefreshToken   string    `json:"refresh_token" gorm:"column:refresh_token;default:null"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" grom:"column:updated_at;autoUpdateTime"`
}

func (u *User) TableName() string {
	return "users"
}
