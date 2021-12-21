package google

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/asumsi/api.inlive/internal/models/auth"
	googleModel "github.com/asumsi/api.inlive/internal/models/google"
	"github.com/asumsi/api.inlive/internal/models/user"
	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Callback(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var (
		oauthConf = &oauth2.Config{
			ClientID:     pkg.GetConfigString(`google_clientId`),
			ClientSecret: pkg.GetConfigString(`google_clientSecret`),
			RedirectURL:  pkg.GetConfigString(`google_redirectURL`),
			Scopes:       []string{pkg.GetConfigString(`google_scopeURL`)},
			Endpoint:     google.Endpoint,
		}
		oauthStateString = pkg.GetConfigString("google_oauth_StateString")
	)

	state := params["state"]
	if state != oauthStateString {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: "Unauthorized", Data: ""})
	}

	code := params["code"]
	if code == "" {
		reason := params["error_reason"]
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: reason, Data: ""})
	}

	token, err := oauthConf.Exchange(context.Background(), code)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	res := googleModel.APIResponse{}
	json.Unmarshal([]byte(response), &res)

	check, err := googleModel.IsUserRegistered(res.Email)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	if !check {
		var userData user.User
		userData.Email = res.Email
		userData.Name = res.Email
		userData.LoginType = "GOOGLE"
		userData.RoleID = int64(1)
		userData.IsActive = true

		user, errUser := auth.Register(&userData)
		if errUser != nil {
			api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: user})
		}
	}

	var req auth.LoginRequest
	req.Email = res.Email
	req.Username = res.Email
	data, err := auth.Authenticate(req, "GOOGLE")
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	accessToken := data.Token
	refreshToken := data.RefreshToken

	result := auth.LoginAttribute{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	api.RespondJSON(w, api.Response{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: result})
}
