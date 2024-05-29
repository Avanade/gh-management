package middleware

import (
	"main/pkg/session"
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func AzureAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session.IsAuthenticated(w, r)
			f(w, r)
		}
	}
}

func GitHubAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session.IsGHAuthenticated(w, r)
			f(w, r)
		}
	}
}

func ManagedIdentityAuth() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session.IsGuidAuthenticated(w, r)
			f(w, r)
		}
	}
}

func IsUserAdmin() Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session.IsUserAdmin(w, r)
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
