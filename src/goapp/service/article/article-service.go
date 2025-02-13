package article

import (
	"main/model"
	"main/repository"
)

type articleService struct {
	Repository *repository.Repository
}

func NewArticleService(repository *repository.Repository) ArticleService {
	return &articleService{repository}
}

func (s *articleService) Insert(article *model.Article) (int64, error) {
	id, err := s.Repository.Article.Insert(article)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *articleService) GetByCategoryId(categoryId int64) ([]model.Article, error) {
	data, err := s.Repository.Article.SelectByCategoryId(categoryId)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (s *articleService) GetById(id int64) (*model.Article, error) {
	data, err := s.Repository.Article.SelectById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *articleService) Update(article *model.Article) error {
	err := s.Repository.Article.Update(article)
	if err != nil {
		return err
	}
	return nil
}
