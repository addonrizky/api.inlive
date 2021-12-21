package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func ExecuteRequest(r *mux.Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	return rr
}

func ExpectResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
