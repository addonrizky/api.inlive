package google

import "github.com/asumsi/api.inlive/internal/models/db"

func countUserByEmail(email string) (total int, err error) {
	sql := db.Connect().
		Table("users").
		Select("count(*) as total").
		Where("email = ?", email).
		Where("login_type = 'GOOGLE'").
		Scan(&total)
	if sql.Error != nil {
		return 0, sql.Error
	}
	return total, nil
}

func IsUserRegistered(email string) (bool, error) {
	count, err := countUserByEmail(email)
	if count < 1 || err != nil {
		return false, err
	}

	return true, err
}
