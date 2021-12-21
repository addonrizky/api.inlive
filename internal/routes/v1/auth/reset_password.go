package auth

import (
	"net/http"

	"github.com/asumsi/api.inlive/internal/models/auth"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	"github.com/gorilla/mux"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {

	var userData auth.ResetPasswordReq

	params := mux.Vars(r)
	key := params["key"]
	emailUser, err := pkg.Decrypter(key)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	userData.Email = emailUser

	err = auth.UpdatePassword(&userData)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	} else {
		api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: ""})
	}

}
