package activitytype

import "main/model"

type ActivityTypeService interface {
	Get() ([]model.ActivityType, error)
}
