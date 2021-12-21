package user

import "time"

type User struct {
	ID              int64      `json:"id" gorm:"column:id"`
	UserName        string     `json:"username"  gorm:"column:username"`
	Password        string     `json:"password"  gorm:"column:password"`
	ConfirmPassword string     `json:"confirm_password" gorm:"-"`
	Name            string     `json:"name"  gorm:"column:name"`
	LoginType       string     `json:"login_type"  gorm:"column:login_type"`
	Email           string     `json:"email"  gorm:"column:email"`
	RoleID          int64      `json:"role_id"  gorm:"column:role_id"`
	IsActive        bool       `json:"is_active"  gorm:"column:is_active"`
	RegisterDate    time.Time  `json:"register_date"  gorm:"column:register_date"`
	UpdatedDate     *time.Time `json:"updated_date"  gorm:"column:updated_date"`
}
