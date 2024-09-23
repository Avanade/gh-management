package activity

import (
	"main/model"
)

type ActivityService interface {
	Get(offset, filter, orderby, ordertype, search, createdBy string) (activities []model.Activity, total int64, err error)
	GetById(id string) (*model.Activity, error)
	Create(activity *model.Activity) (*model.Activity, error)
	Validate(activity *model.Activity) error
}
