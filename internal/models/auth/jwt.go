package auth

import "github.com/dgrijalva/jwt-go"

type (
	JwtClaims struct {
		jwt.StandardClaims
		Name   string `json:"name"`
		Email  string `json:"email"`
		RoleID int64  `json:"role_id"`
	}
	JwtModuleAccess struct {
		ModuleName      string   `json:"module_name"`
		JwtModuleAction []string `json:"action,omitempty"`
	}
)
