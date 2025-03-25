package contributionarea

import "net/http"

type ContributionAreaController interface {
	GetContributionAreas(w http.ResponseWriter, r *http.Request)
	GetContributionAreaById(w http.ResponseWriter, r *http.Request)
	CreateContributionAreas(w http.ResponseWriter, r *http.Request)
	UpdateContributionArea(w http.ResponseWriter, r *http.Request)
}
