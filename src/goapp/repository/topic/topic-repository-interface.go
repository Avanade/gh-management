package topic

import (
	"main/model"
)

type TopicRepository interface {
	GetPopularTopics(opt *model.FilterOptions) ([]model.Topic, error)
	Delete(id int64) error
	Insert(topic string, id int64) error
}
