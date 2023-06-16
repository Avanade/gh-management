package authentication

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/pkg/envvar"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Set up OAuth2 Configurations for GitHub Authentication
func GetGitHubOauthConfig(host string) *oauth2.Config {
	homeUrl := fmt.Sprint(envvar.GetEnvVar("SCHEME", "https"), "://", host)
	oauthConf := &oauth2.Config{
		ClientID:     os.Getenv("GH_CLIENT_ID"),
		ClientSecret: os.Getenv("GH_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		RedirectURL:  homeUrl + "/login/github/callback",
		Endpoint:     githuboauth.Endpoint,
	}
	return oauthConf
}

// Returns GitHub profile of the user
func GetGitHubUserProfile(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}
