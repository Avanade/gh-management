package externallink

import "net/http"

type ExternalLinkController interface {
	GetExternalLinks(w http.ResponseWriter, r *http.Request)
	GetEnabledExternalLinks(w http.ResponseWriter, r *http.Request)
	GetExternalLinkById(w http.ResponseWriter, r *http.Request)
	CreateExternalLink(w http.ResponseWriter, r *http.Request)
	UpdateExternalLinkById(w http.ResponseWriter, r *http.Request)
	RemoveExternalLinkById(w http.ResponseWriter, r *http.Request)
}
