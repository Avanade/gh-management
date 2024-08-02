package middleware

import (
	"main/pkg/session"
	"net/http"
)

func IsUserAdmin() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if session.IsUserAdminMW(w, r) {
				f(w, r)
			}
		}
	}
}
