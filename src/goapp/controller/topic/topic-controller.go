package topic

import (
	"encoding/json"
	"main/model"
	"main/pkg/appinsights_wrapper"
	"main/service"
	"net/http"
	"strconv"
)

type topicController struct {
	*service.Service
}

func NewTopicController(serv *service.Service) TopicController {
	return &topicController{
		Service: serv,
	}
}

func (c *topicController) GetPopularTopics(w http.ResponseWriter, r *http.Request) {
	logger := appinsights_wrapper.NewClient()
	defer logger.EndOperation()

	var topics []model.Topic
	var err error

	params := r.URL.Query()
	if params.Has("offset") && params.Has("rowCount") {
		filter, _ := strconv.Atoi(params["rowCount"][0])
		offset, _ := strconv.Atoi(params["offset"][0])
		opt := model.FilterOptions{
			Filter: filter,
			Offset: offset,
		}
		topics, err = c.Topic.Get(&opt)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		topics, err = c.Topic.Get(nil)
		if err != nil {
			logger.TrackException(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(topics)
}
