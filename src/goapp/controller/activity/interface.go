package activity

import "net/http"

type ActivityController interface {
	GetActivities(w http.ResponseWriter, r *http.Request)
	GetActivityById(w http.ResponseWriter, r *http.Request)
	CreateActivity(w http.ResponseWriter, r *http.Request)
}
