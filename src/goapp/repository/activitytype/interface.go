package activitytype

import "main/model"

type ActivityTypeRepository interface {
	Insert(body *model.ActivityType) (*model.ActivityType, error)
	Select() ([]model.ActivityType, error)
}
