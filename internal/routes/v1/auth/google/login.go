package google

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/asumsi/api.inlive/pkg"
	"github.com/asumsi/api.inlive/pkg/api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Login(w http.ResponseWriter, r *http.Request) {

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

	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		api.RespondJSON(w, api.Response{Code: http.StatusBadRequest, Message: err.Error(), Data: ""})
	}

	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	resp := api.Response{Code: http.StatusOK, Message: http.StatusText(http.StatusOK), Data: url}
	api.RespondJSON(w, resp)
}
