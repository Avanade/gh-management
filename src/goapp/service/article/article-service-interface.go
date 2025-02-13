package article

import (
	"main/model"
)

type ArticleService interface {
	Insert(article *model.Article) (int64, error)
	GetByCategoryId(id int64) ([]model.Article, error)
	GetById(id int64) (*model.Article, error)
	Update(article *model.Article) error
}
