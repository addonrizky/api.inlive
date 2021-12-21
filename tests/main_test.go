package test

import (
	"os"
	"testing"

	"github.com/asumsi/api.inlive/internal/models/db"
	"github.com/asumsi/api.inlive/internal/models/stream"
	"github.com/asumsi/api.inlive/internal/models/user"
	"github.com/asumsi/api.inlive/internal/routes"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Application struct {
	DB     *gorm.DB
	router *mux.Router
}

var App Application

func init() {
	App.DB = db.Connect()
	App.router = mux.NewRouter()
	routes.InitRoutes(App.router)
}

// func startHttpServer() {
// 	// Bind to a port and pass our router in
// 	port := pkg.GetConfigString("PORT")
// 	if port != "" {
// 		log.Fatal(http.ListenAndServe(":"+port, A.router))
// 	} else {
// 		log.Fatal(http.ListenAndServe(":8080", A.router))
// 	}
// }
func TestMain(m *testing.M) {
	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	App.DB.AutoMigrate(&user.User{}, &stream.Stream{})
	App.DB.Migrator().HasTable("users")
	App.DB.Migrator().HasTable("streams")
}

func clearTable() {
	App.DB.Migrator().DropTable(&user.User{}, &stream.Stream{})
}
