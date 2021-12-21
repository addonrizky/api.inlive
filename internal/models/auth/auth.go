package auth

import (
	"errors"
	"time"

	"github.com/asumsi/api.inlive/internal/models/db"
	"github.com/asumsi/api.inlive/internal/models/user"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/dgrijalva/jwt-go"
)

func Register(newUser *user.User) (user.User, error) {

	if newUser.Password != newUser.ConfirmPassword {
		return user.User{}, errors.New("password is not same")
	}

	hashPass, err := pkg.HashPassword(newUser.Password)
	if err != nil {
		return user.User{}, err
	}

	newUser.Password = hashPass
	return user.CreateUser(newUser)
}

func Authenticate(req LoginRequest, loginType string) (resp AuthResponse, err error) {
	var authUser user.User

	authUser.Email = req.Email
	authUser.UserName = req.Email
	authUser.RoleID = int64(1)

	if loginType == "APP" {
		if pkg.IsEmailValid(req.Email) {
			authUser, err = user.GetUserByEmail(req.Email)
		} else {
			authUser, err = user.GetUserByUsername(req.Username)
		}
		if authUser.ID == 0 || err != nil {
			return resp, errors.New("Username/Email tidak terdaftar")
		}
		if !authUser.IsActive {
			return resp, errors.New("Username/Email tidak aktif")
		}

		if !pkg.CheckPasswordHash(req.Password, authUser.Password) {
			return resp, errors.New("password salah")
		}
	}

	resultsAll, err := generateToken(authUser.UserName, authUser.Email, authUser.RoleID)
	return resultsAll, err
}

func generateToken(name, email string, roleID int64) (resp AuthResponse, err error) {
	appName := pkg.GetConfigString(`app_name`)
	jwtMethod := pkg.GetConfigString(`jwt_method`)
	jwtSecret := pkg.GetConfigString(`jwt_secret`)
	jwtLifespan := pkg.GetConfigDuration(`jwt_lifespan`)
	jwtLifespanRefresh := pkg.GetConfigDuration(`jwt_refresh_lifespan`)
	resp.Token, err = pkg.GenerateJwtToken(JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    appName,
			ExpiresAt: time.Now().Add(jwtLifespan).Unix(),
		},
		Name:   name,
		Email:  email,
		RoleID: roleID,
	}, jwtMethod, jwtSecret)
	resp.RefreshToken, err = pkg.GenerateJwtToken(JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    appName,
			ExpiresAt: time.Now().Add(jwtLifespanRefresh).Unix(),
		},
		Name:  name,
		Email: email,
	}, jwtMethod, jwtSecret)
	return resp, nil
}

func UpdatePassword(req *ResetPasswordReq) error {
	if req.Password != req.ConfirmPassword {
		return errors.New("password is not same")
	}

	hashPass, err := pkg.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashPass

	err = UpdatePasswordByEmail(req)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePasswordByEmail(req *ResetPasswordReq) error {
	now := time.Now()
	user := &user.User{Email: req.Email, Password: req.Password, UpdatedDate: &now}

	sql := db.Connect().Table("users").
		Where("email = ?", user.Email).
		Where("is_active = ?", true).
		Where("login_type = ?", "APP").
		Updates(&user)
	if sql.Error != nil {
		return sql.Error
	}

	return nil
}
