package middleware

import (
	"main/pkg/session"
	"net/http"
)

func GuidAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if session.IsGuidAuthenticated(w, r) {
				f(w, r)
			}
		}
	}
}
