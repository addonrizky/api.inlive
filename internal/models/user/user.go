package user

import (
	"time"

	"github.com/asumsi/api.inlive/internal/models/db"
)

func CreateUser(req *User) (User, error) {
	user := User{UserName: req.UserName, Password: req.Password, Name: req.Name, LoginType: req.LoginType, Email: req.Email, RoleID: req.RoleID, IsActive: req.IsActive, RegisterDate: time.Now()}

	sql := db.Connect().Table("users").Create(&user)

	return user, sql.Error
}

func GetUserByEmail(email string) (res User, err error) {
	sql := db.Connect().Table("users").Select("id, username, name, password, login_type, is_active").Where("email = ?", email).Scan(&res)
	if sql.Error != nil {
		return res, sql.Error
	}
	return res, nil
}

func GetUserByUsername(username string) (res User, err error) {
	sql := db.Connect().Table("users").Select("id, username, name, password, login_type, is_active").Where("username = ?", username).Scan(&res)
	if sql.Error != nil {
		return res, sql.Error
	}
	return res, nil
}
