package activity

import "main/model"

type ActivityRepository interface {
	Insert(activity *model.Activity) (*model.Activity, error)
	Select() ([]model.Activity, error)
	SelectByOptions(offset, filter int64, orderBy, orderType, search, createdBy string) ([]model.Activity, error)
	SelectById(id int64) (*model.Activity, error)
	Total() (int64, error)
	TotalByOptions(search, createdBy string) (int64, error)
}
