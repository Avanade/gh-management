package routes

import (
	auth "main/pkg/authentication"
	session "main/pkg/session"
	"net/http"
)

func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {

	// Check session
	session, err := session.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	// Get GitHub temporary code
	code := r.URL.Query().Get("code")

	ghAccessToken := auth.GetGithubAccessToken(code)
	ghProfile := auth.GetGitHubUserProfile(ghAccessToken)

	// Save GitHub auth data on session cookies
	session.Values["ghAccessToken"] = ghAccessToken
	session.Values["ghProfile"] = ghProfile

	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
