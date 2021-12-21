package pkg

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJwtToken(claims jwt.Claims, method, secret string) (token string, err error) {
	jwtToken := jwt.NewWithClaims(
		jwt.GetSigningMethod(method),
		claims,
	)
	return jwtToken.SignedString([]byte(secret))
}

func ValidateJwtToken(tokenString, method, secret string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(method) != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	return jwtToken, err
}
