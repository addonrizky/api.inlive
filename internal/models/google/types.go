package google

type APIResponse struct {
	Sub           string `json:"sub"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	HD            string `json:"hd"`
}
