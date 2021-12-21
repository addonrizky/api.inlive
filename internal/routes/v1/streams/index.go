package streams

import (
	"net/http"

	"github.com/asumsi/api.inlive/pkg/router"
	"github.com/pion/webrtc/v3"
)

type Controller struct {
	Sessions map[string]*webrtc.PeerConnection
	Ports    map[uint16]bool
	FFmpegs  map[string]int
}

func GetRoutesHandlers() []router.Route {
	controller := &Controller{
		Sessions: make(map[string]*webrtc.PeerConnection, 34000),
		Ports:    make(map[uint16]bool, 17000),
		FFmpegs:  make(map[string]int, 17000),
	}
	routes := []router.Route{
		{Path: "/streams", Handler: List, Method: http.MethodGet},
		{Path: "/streams/create", Handler: Create, Method: http.MethodPost},
		{Path: "/streams/{id}/init", Handler: controller.Init, Method: http.MethodPost},
		{Path: "/streams/{id}/start", Handler: controller.Start, Method: http.MethodPost},
		{Path: "/streams/{id}/end", Handler: controller.End, Method: http.MethodPost},
		{Path: "/streams/{id}", Handler: controller.Get, Method: http.MethodGet},
		{Path: "/streams/{id}", Handler: Update, Method: http.MethodPut},
		{Path: "/streams/{id}", Handler: Delete, Method: http.MethodDelete},
	}
	return routes
}
