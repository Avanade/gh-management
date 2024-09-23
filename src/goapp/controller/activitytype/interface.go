package activitytype

import (
	"net/http"
)

type ActivityTypeController interface {
	GetActivityTypes(w http.ResponseWriter, r *http.Request)
}
