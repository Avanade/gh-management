package category

import (
	"main/model"
	"main/repository"
)

type categoryService struct {
	Repository *repository.Repository
}

func NewCategoryService(repository *repository.Repository) CategoryService {
	return &categoryService{repository}
}

func (s *categoryService) Insert(category *model.Category) (int64, error) {
	id, err := s.Repository.Category.Insert(category)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *categoryService) GetAll() ([]model.Category, error) {
	return s.Repository.Category.Select()
}

func (s *categoryService) GetById(id int64) (*model.Category, error) {
	data, err := s.Repository.Category.SelectById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *categoryService) Update(category *model.Category) error {
	err := s.Repository.Category.Update(category)
	if err != nil {
		return err
	}
	return nil
}
