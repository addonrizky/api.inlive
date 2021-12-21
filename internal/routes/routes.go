package routes

import (
	"fmt"
	"net/http"

	"github.com/asumsi/api.inlive/internal/routes/v1/streams"
	"github.com/asumsi/api.inlive/pkg/api"
	"github.com/gorilla/mux"
)

var API_PREFIX = "/v1"

func registerRoute(r *mux.Router, path string, handler api.HttpHandler, method string) {
	fullpath := API_PREFIX + path
	fmt.Println(`Registered route: ` + method + " " + fullpath)
	r.HandleFunc(fullpath, handler).Methods(method, http.MethodOptions)
}

func InitRoutes(r *mux.Router) {
	// IMPORTANT: you must specify an OPTIONS method matcher for the middleware to set CORS headers
	streams := streams.GetRoutesHandlers()
	for _, route := range streams {
		registerRoute(r, route.Path, route.Handler, route.Method)
	}

	// auths := auth.GetRoutesHandlers()
	// for _, route := range auths {
	// 	registerRoute(r, route.Path, route.Handler, route.Method)
	// }

	r.HandleFunc(`/`, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`ok`))
	})
}
