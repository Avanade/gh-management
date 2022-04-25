package session

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	auth "webserver/pkg/authentication"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

var (
	Store *sessions.FilesystemStore
)

func InitializeSession() {
	Store = sessions.NewFilesystemStore(os.TempDir(), []byte("secret"))
	Store.MaxLength(math.MaxInt64)
	gob.Register(map[string]interface{}{})
}

func IsAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	// Check session if there is saved user profile
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := session.Values["profile"]; !ok {

		// Asks user to login if there is no saved user profile
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

	} else {

		// If there is a user profile saved
		authenticator, err := auth.NewAuthenticator()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve current token data
		refreshToken := fmt.Sprintf("%s", session.Values["refresh_token"])
		expiry, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s", session.Values["expiry"]))
		today := time.Now().UTC()

		if today.After(expiry) {
			// Attempt to refresh token if access token is already expired
			ts := authenticator.Config.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken, Expiry: expiry})
			newToken, err := ts.Token()
			if err != nil {
				fmt.Printf("ERROR REFRESHING TOKEN\n")
				if rErr, ok := err.(*oauth2.RetrieveError); ok {
					details := new(ErrorDetails)
					if err := json.Unmarshal(rErr.Body, details); err != nil {
						panic(err)
					}

					fmt.Println(details.Error, details.ErrorDescription)
				}

				// Delete current session
				session.Options.MaxAge = -1
				session.Save(r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// Ask user to re-login
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)

			} else if newToken != nil {
				fmt.Printf("TOKEN REFRESHED\n")

				// Save new token data
				session.Values["refresh_token"] = newToken.RefreshToken
				session.Values["expiry"] = newToken.Expiry.String()
				session.Values["access_token"] = newToken.AccessToken
				err = session.Save(r, w)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		next(w, r)
	}
}

type ErrorDetails struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
