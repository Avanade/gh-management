package topic

import "net/http"

type TopicController interface {
	GetPopularTopics(w http.ResponseWriter, r *http.Request)
}
