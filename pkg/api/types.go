package api

import "net/http"

type Pagination struct {
	Page        int `json:"page"`
	DataPerPage int `json:"dataPerPage"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type HttpHandler func(http.ResponseWriter, *http.Request)
