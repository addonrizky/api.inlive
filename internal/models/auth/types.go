package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginAttribute struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Result  LoginAttribute `json:"result"`
}

type ResetPasswordReq struct {
	Email           string `json:"email" `
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type EmailResetPasswordReq struct {
	Email string `json:"email" `
}

type EmailResetPasswordResponse struct {
	Code              int    `json:"code"`
	Message           string `json:"message"`
	ResetPasswordLink string `json:"reset_password_link"`
}
