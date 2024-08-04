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
	FullName     string `json:"name" validate:"required,min=3,max=100"`
	PhotoProfile string `json:"picture" validate:"required,min=3,max=500"`
	RefreshToken string `json:"refresh_token" validate:"required,min=50,max=500"`
}

type AddRefreshTokenReq struct {
	Email        string `json:"email" validate:"required,email,min=5,max=100"`
	RefreshToken string `json:"refresh_token" validate:"required,min=50,max=500"`
}
