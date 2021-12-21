package db

import (
	"fmt"

	"github.com/asumsi/api.inlive/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func init() {
	dbHost := pkg.GetConfigString(`DB_HOST`)
	dbUser := pkg.GetConfigString(`DB_USER`)
	dbPass := pkg.GetConfigString(`DB_PASS`)
	dbName := pkg.GetConfigString(`DB_NAME`)
	debug := pkg.GetConfigBool(`debug`)
	logLevel := logger.Error
	if debug{
		logLevel = logger.Info
	}

	connection := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)
	fmt.Println(connection)
	var err error
	Db, err = gorm.Open(postgres.Open(connection), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil && pkg.GetConfigBool("debug") {
		fmt.Println(err)
	}
}

func Connect() *gorm.DB {
	return Db
}
