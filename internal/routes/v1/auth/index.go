package auth

import (
	"net/http"

	"github.com/asumsi/api.inlive/internal/routes/v1/auth/google"
	"github.com/asumsi/api.inlive/pkg/router"
)

func GetRoutesHandlers() []router.Route {
	routes := []router.Route{
		{Path: "/auth/google/callback", Handler: google.Callback, Method: http.MethodGet},
		{Path: "/auth/google/login", Handler: google.Login, Method: http.MethodPost},
	}
	return routes
}
