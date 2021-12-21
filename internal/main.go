package main

import (
	"fmt"
	"log"
	"net/http"

	//"github.com/asumsi/api.inlive/internal/models/db/migration"
	"github.com/asumsi/api.inlive/internal/routes"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	_ "github.com/asumsi/api.inlive/docs"
)

type Application struct {
	DB     *gorm.DB
	router *mux.Router
}

var App Application

func init() {
	if pkg.GetConfigBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
	pkg.ValidatorInit()
}

// @title Livestream API
// @version 1.0
// @description Livestream API enable you to live stream.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host api.inlive.app
// @BasePath /v1
func main() {
	//migration.Migrate();
	r := mux.NewRouter()
	routes.InitRoutes(r)

	cors := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		// handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin"}),
		// handlers.AllowCredentials(),
	)

	// Bind to a port and pass our router in
	port := pkg.GetConfigString("PORT")
	if port != "" {
		log.Fatal(http.ListenAndServe(":"+port, cors(r)))
	} else {
		log.Fatal(http.ListenAndServe(":8080", cors(r)))
	}

}
