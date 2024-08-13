package dto

type CreateReq struct {
	UserId   string `json:"user_id" validate:"required,min=21,max=21"`
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	FullName string `json:"full_name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type UpsertReq struct {
	UserId       string `json:"user_id" validate:"required,min=21,max=21"`
	Email        string `json:"email" validate:"required,email,min=5,max=100"`
	FullName     string `json:"full_name" validate:"required,min=3,max=100"`
	PhotoProfile string `json:"photo_profile" validate:"required,min=3,max=500"`
	RefreshToken string `json:"refresh_token" validate:"required,min=50,max=500"`
}

type UpdateProfileReq struct {
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	FullName string `json:"full_name" validate:"omitempty,min=3,max=100"`
	Whatsapp string `json:"whatsapp" validate:"omitempty,min=10,max=20"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type UpdatePasswordReq struct {
	Email       string `json:"email" validate:"required,email,min=5,max=100"`
	Password    string `json:"password" validate:"required,min=5,max=100"`
	NewPassword string `json:"new_password" validate:"required,min=5,max=100"`
}

type UpdateEmailReq struct {
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	NewEmail string `json:"new_email" validate:"required,email,min=5,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100"`
}

type UpdatePhotoProfileReq struct {
	Email          string `json:"email" validate:"required,email,min=5,max=100"`
	PhotoProfileId string `json:"new_photo_profile_id" validate:"required,min=10,max=100"`
	PhotoProfile   string `json:"new_photo_profile" validate:"required,min=10,max=500"`
}

type VerifyUpdateEmailReq struct {
	Email    string `json:"email" validate:"required,email,min=5,max=100"`
	NewEmail string `json:"new_email" validate:"required,email,min=5,max=100"`
	Otp      string `json:"otp" validate:"required,max=6"`
}

type VerifyOtpReq struct {
	Email string `json:"email" validate:"required,email,min=5,max=100"`
	Otp   string `json:"otp" validate:"required,max=6"`
}

type AddRefreshTokenReq struct {
	Email        string `json:"email" validate:"required,email,min=5,max=100"`
	RefreshToken string `json:"refresh_token" validate:"required,min=50,max=500"`
}
