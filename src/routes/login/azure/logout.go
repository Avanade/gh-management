package routes

import (
	session "main/pkg/session"
	"net/http"
	"net/url"
	"os"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutUrl, err := url.Parse("https://login.microsoftonline.com/" + os.Getenv("tenantid") + "/oauth2/logout?client_id=" + os.Getenv("clientid") + "&post_logout_redirect_uri=" + os.Getenv("homeurl"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
