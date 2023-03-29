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

	"main/models"
	auth "main/pkg/authentication"

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
	url := fmt.Sprintf("/loginredirect?redirect=%v", r.URL)
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		c := http.Cookie{
			Name:   "auth-session",
			MaxAge: -1}
		http.SetCookie(w, &c)
		cgh := http.Cookie{
			Name:   "gh-auth-session",
			MaxAge: -1}
		http.SetCookie(w, &cgh)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	if _, ok := session.Values["profile"]; !ok {

		// Asks user to login if there is no saved user profile
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)

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
				// fmt.Printf("ERROR REFRESHING TOKEN\n")
				if rErr, ok := err.(*oauth2.RetrieveError); ok {
					details := new(ErrorDetails)
					if err := json.Unmarshal(rErr.Body, details); err != nil {
						panic(err)
					}

					fmt.Println(details.Error, details.ErrorDescription)
				}

				// Log out the user if the attempt to refresh the token failed
				http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)

			} else if newToken != nil {
				// fmt.Printf("TOKEN REFRESHED\n")

				// Save new token data
				session.Values["refresh_token"] = newToken.RefreshToken
				session.Values["expiry"] = newToken.Expiry.UTC().Format("2006-01-02 15:04:05")
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

func IsGHAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Check session if there is saved user profile
	url := fmt.Sprintf("/gitredirect?redirect=%v", r.URL)
	session, err := Store.Get(r, "gh-auth-session")
	if err != nil {
		c := http.Cookie{
			Name:   "gh-auth-session",
			MaxAge: -1}
		http.SetCookie(w, &c)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}

	if _, ok := session.Values["ghProfile"]; !ok || !session.Values["ghIsValid"].(bool) {

		// Asks user to login if there is no saved user profile
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	} else {
		next(w, r)
	}
}

func IsGuidAuthenticated(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// Check header if authenticated
	_, err := auth.VerifyAccessToken(r)
	// RETURN ERROR
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// RETURN SUCCESS
	next(w, r)
}

func GetGitHubUserData(w http.ResponseWriter, r *http.Request) (models.TypGitHubUser, error) {
	session, err := Store.Get(r, "gh-auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return models.TypGitHubUser{LoggedIn: false}, err
	}
	var gitHubUser models.TypGitHubUser

	if _, ok := session.Values["ghProfile"]; ok {
		err = json.Unmarshal([]byte(fmt.Sprintf("%s", session.Values["ghProfile"])), &gitHubUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return models.TypGitHubUser{LoggedIn: false}, err
		}

		if _, okIsValid := session.Values["ghIsValid"]; okIsValid {
			gitHubUser.IsValid = session.Values["ghIsValid"].(bool)
		}

		if _, ok := session.Values["ghIsDirect"]; ok {
			gitHubUser.IsDirect = session.Values["ghIsDirect"].(bool)
		}
		if _, ok := session.Values["ghIsEnterpriseMember"]; ok {
			gitHubUser.IsEnterpriseMember = session.Values["ghIsEnterpriseMember"].(bool)
		}

		gitHubUser.AccessToken = fmt.Sprintf("%s", session.Values["ghAccessToken"])
		gitHubUser.LoggedIn = true
	} else {
		gitHubUser.LoggedIn = false
	}

	return gitHubUser, nil
}

func GetState(w http.ResponseWriter, r *http.Request) (string, error) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	if _, ok := session.Values["state"]; ok {

		return fmt.Sprintf("%s", session.Values["state"]), nil

	}
	return "", nil
}

func RemoveGitHubAccount(w http.ResponseWriter, r *http.Request) error {
	session, err := Store.Get(r, "gh-auth-session")
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)

	if err != nil {
		return err
	}

	return nil
}

func IsUserAdminMW(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	isAdmin, _ := IsUserAdmin(w, r)
	if isAdmin {
		next(w, r)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Check if user is an admin
func IsUserAdmin(w http.ResponseWriter, r *http.Request) (bool, error) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		return false, err
	}

	isUserAdmin := false
	if session.Values["isUserAdmin"] != nil {
		isUserAdmin = session.Values["isUserAdmin"].(bool)
	}
	return isUserAdmin, nil
}

// Get user photo
func GetUserPhoto(w http.ResponseWriter, r *http.Request) (bool, string, error) {
	session, err := Store.Get(r, "auth-session")
	if err != nil {
		return false, "", err
	}

	if session.Values["userHasPhoto"] != nil {
		return session.Values["userHasPhoto"].(bool), fmt.Sprintf("%s", session.Values["userPhoto"]), nil
	} else {
		return false, "", nil
	}
}

type ErrorDetails struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
