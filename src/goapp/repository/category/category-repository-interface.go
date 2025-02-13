package category

import (
	"main/model"
)

type CategoryRepository interface {
	Insert(category *model.Category) (int64, error)
	Select() ([]model.Category, error)
	SelectById(categoryid int64) (*model.Category, error)
	Update(category *model.Category) error
}
