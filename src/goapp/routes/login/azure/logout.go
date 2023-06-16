package routes

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"main/pkg/envvar"
	"main/pkg/session"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	azSession, err := session.Store.Get(r, "auth-session")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	azSession.Options.MaxAge = -1
	err = azSession.Save(r, w)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := session.GetGitHubUserData(w, r)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.LoggedIn {
		session.RemoveGitHubAccount(w, r)
	}

	homeUrl := fmt.Sprint(envvar.GetEnvVar("SCHEME", "https"), "://", r.Host)
	logoutUrl, err := url.Parse("https://login.microsoftonline.com/" + os.Getenv("TENANT_ID") + "/oauth2/logout?client_id=" + os.Getenv("CLIENT_ID") + "&post_logout_redirect_uri=" + homeUrl)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
