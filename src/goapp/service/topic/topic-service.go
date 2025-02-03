package topic

import (
	"errors"
	"main/model"
	"main/repository"
	"strconv"
)

type topicService struct {
	Repository *repository.Repository
}

func NewTopicService(repository *repository.Repository) TopicService {
	return &topicService{repository}
}

func (s *topicService) Get(opt *model.FilterOptions) ([]model.Topic, error) {
	var topics []model.Topic

	if opt == nil {
		return nil, errors.New("FilterOptions is nil")
	}

	data, err := s.Repository.Topic.GetPopularTopics(opt)
	if err != nil {
		return nil, err
	}
	topics = data
	return topics, nil
}

func (s *topicService) Delete(id string) error {
	parsedId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	return s.Repository.Topic.Delete(parsedId)
}

func (s *topicService) Insert(topic string, id int64) error {
	err := s.Repository.Topic.Insert(topic, id)
	if err != nil {
		return err
	}
	return nil
}
