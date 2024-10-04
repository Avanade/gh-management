package osscontributionsponsor

import "net/http"

type OSSContributionSponsorController interface {
	GetOssContributionSponsors(w http.ResponseWriter, r *http.Request)
	GetEnabledOssContributionSponsors(w http.ResponseWriter, r *http.Request)
	CreateOssContributionSponsor(w http.ResponseWriter, r *http.Request)
	UpdateOssContributionSponsor(w http.ResponseWriter, r *http.Request)
}
