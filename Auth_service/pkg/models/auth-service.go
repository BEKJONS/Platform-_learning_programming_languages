package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       string `json:"user_id"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}
type ResetPassReq struct {
	Email    string `json:"email"`
	Password string `json:"new_password"`
	Code     string `json:"code"`
}
type UpdatePasswordReq struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

type Message struct {
	Message string `json:"message"`
}
type AcceptCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
