package article

import (
	"main/model"
)

type ArticleRepository interface {
	Insert(article *model.Article) (int64, error)
	SelectByCategoryId(categoryId int64) ([]model.Article, error)
	SelectById(id int64) (*model.Article, error)
	Update(category *model.Article) error
}
