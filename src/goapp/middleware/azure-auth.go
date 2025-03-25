package middleware

import (
	"main/pkg/session"
	"net/http"
)

func AzureAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if session.IsAuthenticated(w, r) {
				f(w, r)
			}
		}
	}
}
