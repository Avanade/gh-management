package middleware

import (
	"main/pkg/session"
	"net/http"
)

func GitHubAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if session.IsGHAuthenticated(w, r) {
				f(w, r)
			}
		}
	}
}
