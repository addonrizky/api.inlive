package auth

import (
	"net/http"

	"github.com/asumsi/api.inlive/internal/models/auth"
	"github.com/asumsi/api.inlive/pkg/api"
)

func Login(w http.ResponseWriter, r *http.Request) {

	var req auth.LoginRequest

	data, err := auth.Authenticate(req, "APP")

	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	result := auth.LoginAttribute{
		Token:        data.Token,
		RefreshToken: data.RefreshToken,
	}
	api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusCreated), Data: result})

}
