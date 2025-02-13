package category

import (
	"main/model"
)

type CategoryService interface {
	Insert(category *model.Category) (int64, error)
	GetAll() ([]model.Category, error)
	GetById(id int64) (*model.Category, error)
	Update(category *model.Category) error
}
