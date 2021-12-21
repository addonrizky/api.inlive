package auth

import (
	"net/http"

	"github.com/asumsi/api.inlive/internal/models/auth"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
)

func SendEmailResetPassword(w http.ResponseWriter, r *http.Request) {

	var req auth.ResetPasswordReq

	encryptEmail, err := pkg.Encrypter(req.Email)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	// send email via SMTP
	type resetPassword struct {
		Url string
	}
	prefixLink := pkg.GetConfigString("url")
	resetPasswordUrl := prefixLink + "/api/v1/reset-password?key=" + encryptEmail
	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: resetPassword{Url: resetPasswordUrl}})

}
