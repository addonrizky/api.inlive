package router

import "github.com/asumsi/api.inlive/pkg/api"

type Route struct {
	Path    string
	Handler api.HttpHandler
	Method  string
}
