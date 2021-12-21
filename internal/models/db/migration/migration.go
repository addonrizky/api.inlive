package migration

import (
	"github.com/asumsi/api.inlive/pkg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate() {
	username := pkg.GetConfigString("db_user")
	password := pkg.GetConfigString("db_pass")
	host := pkg.GetConfigString("db_host")
	port := pkg.GetConfigString("db_port")
	db_name := pkg.GetConfigString("db_name")
    m, err := migrate.New(
        "file://internal/models/db/migration/migrations",
        "postgres://"+username+":"+password+"@"+host+":"+port+"/"+db_name+"?sslmode=disable")

	if err != nil {
		panic("Migration error: problem occured when creating db instance:\n"+err.Error())
	}

    err = m.Up()
    if err != nil && err!=migrate.ErrNoChange {
	    panic("Migration error: problem occured when migrating:\n"+err.Error())
    }
}