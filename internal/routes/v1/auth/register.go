package auth

import (
	"net/http"

	"github.com/asumsi/api.inlive/internal/models/auth"
	"github.com/asumsi/api.inlive/internal/models/user"
	"github.com/asumsi/api.inlive/pkg/api"
)

func Register(w http.ResponseWriter, r *http.Request) {

	var userData user.User

	userData.IsActive = true

	user, err := auth.Register(&userData)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	} else {
		api.RespondJSON(w, api.Response{Code: http.StatusCreated, Message: http.StatusText(http.StatusCreated), Data: user})
	}
}
