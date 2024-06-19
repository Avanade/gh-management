package main

import (
	"net/http"
	"strconv"
	"time"

	ev "main/pkg/envvar"
	rtApi "main/routes/api"
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

func checkFailedApprovalRequests() {
	// TIMER SERVICE
	freq := ev.GetEnvVar("APPROVALREQUESTS_RETRY_FREQ", "15")
	freqInt, _ := strconv.ParseInt(freq, 0, 64)
	if freq > "0" {
		for range time.NewTicker(time.Duration(freqInt) * time.Minute).C {
			go rtApi.ReprocessRequestApproval()

			go rtApi.ReprocessCommunityApprovalRequestCommunities()
			go rtApi.ReprocessCommunityApprovalRequestOrganizationAccess()
			go rtApi.ReprocessCommunityApprovalRequestGitHubCoPilots()
			go rtApi.ReprocessCommunityApprovalRequestNewOrganizations()
		}
	}
}
