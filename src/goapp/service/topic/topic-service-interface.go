package topic

import (
	"main/model"
)

type TopicService interface {
	Get(opt *model.FilterOptions) ([]model.Topic, error)
	Delete(id string) error
	Insert(topic string, id int64) error
}
