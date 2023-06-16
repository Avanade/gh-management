package main

import (
	"net/http"
	"strconv"
	"time"

	ev "main/pkg/envvar"
	"main/pkg/session"
	rtApi "main/routes/api"

	"github.com/codegangsta/negroni"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.next.ServeHTTP(w, r)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// Verifies authentication before loading the page.
func loadAzAuthPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func loadGuidAuthApi(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsGuidAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func loadAzGHAuthPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.HandlerFunc(session.IsGHAuthenticated),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func loadAdminPage(f func(w http.ResponseWriter, r *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(session.IsAuthenticated),
		negroni.HandlerFunc(session.IsUserAdminMW),
		negroni.Wrap(http.HandlerFunc(f)),
	)
}

func checkFailedApprovalRequests() {
	// TIMER SERVICE
	freq := ev.GetEnvVar("APPROVALREQUESTS_RETRY_FREQ", "15")
	freqInt, _ := strconv.ParseInt(freq, 0, 64)
	if freq > "0" {
		for range time.NewTicker(time.Duration(freqInt) * time.Minute).C {
			go rtApi.ReprocessRequestApproval()
			go rtApi.ReprocessRequestCommunityApproval()
		}
	}
}
