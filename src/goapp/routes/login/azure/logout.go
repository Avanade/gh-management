package routes

import (
	session "main/pkg/session"
	"net/http"
	"net/url"
	"os"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	azSession, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	azSession.Options.MaxAge = -1
	err = azSession.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := session.GetGitHubUserData(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user.LoggedIn {
		session.RemoveGitHubAccount(w, r)
	}

	logoutUrl, err := url.Parse("https://login.microsoftonline.com/" + os.Getenv("TENANT_ID") + "/oauth2/logout?client_id=" + os.Getenv("CLIENT_ID") + "&post_logout_redirect_uri=" + os.Getenv("HOME_URL"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
