package article

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
)

type articleRepository struct {
	*db.Database
}

func NewArticleRepository(db *db.Database) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) Insert(article *model.Article) (int64, error) {
	row, err := r.QueryRow("usp_GuidanceCategoryArticle_Insert",
		sql.Named("Name", article.Name),
		sql.Named("Url", article.Url),
		sql.Named("Body", article.Body),
		sql.Named("CategoryId", article.CategoryId),
		sql.Named("CreatedBy", article.CreatedBy),
		sql.Named("ModifiedBy", article.ModifiedBy),
	)

	if err != nil {
		return 0, err
	}

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *articleRepository) SelectByCategoryId(id int64) ([]model.Article, error) {

	rows, err := r.Query("[dbo].[usp_GuidanceCategoryArticle_Select_ByCategoryID]",
		sql.Named("Id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	articles := []model.Article{}

	for _, v := range mapRows {
		var article model.Article
		article.Id = v["Id"].(int64)
		article.Name = v["Name"].(string)
		article.Url = v["Url"].(string)
		article.Body = v["Body"].(string)
		article.CategoryId = v["GuidanceCategoryId"].(int64)
		if v["Created"] != nil {
			article.Created = v["Created"].(time.Time)
		}
		article.CreatedBy = v["CreatedBy"].(string)
		if v["Modified"] != nil {
			article.Modified = v["Modified"].(time.Time)
		}
		if v["ModifiedBy"] != nil {
			article.ModifiedBy = v["ModifiedBy"].(string)
		}
		article.Category.Id = v["GuidanceCategoryId"].(int64)
		article.Category.Name = v["CategoryName"].(string)

		articles = append(articles, article)
	}

	return articles, nil
}

func (r *articleRepository) SelectById(id int64) (*model.Article, error) {
	var article model.Article
	rows, err := r.Query("[dbo].[usp_GuidanceCategoryArticle_Select_ByID]",
		sql.Named("Id", id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		article.Id = v["Id"].(int64)
		article.Name = v["Name"].(string)
		article.Url = v["Url"].(string)
		article.Body = v["Body"].(string)
		article.CategoryId = v["GuidanceCategoryId"].(int64)
		if v["Created"] != nil {
			article.Created = v["Created"].(time.Time)
		}
		article.CreatedBy = v["CreatedBy"].(string)
		if v["Modified"] != nil {
			article.Modified = v["Modified"].(time.Time)
		}
		if v["ModifiedBy"] != nil {
			article.ModifiedBy = v["ModifiedBy"].(string)
		}

		article.Category.Id = v["GuidanceCategoryId"].(int64)
		article.Category.Name = v["CategoryName"].(string)

	}

	return &article, nil
}

func (r *articleRepository) Update(article *model.Article) error {
	row, err := r.QueryRow("[dbo].[usp_GuidanceCategoryArticle_Update]",
		sql.Named("Id", article.Id),
		sql.Named("Name", article.Name),
		sql.Named("Url", article.Url),
		sql.Named("Body", article.Body),
		sql.Named("CategoryId", article.CategoryId),
		sql.Named("CreatedBy", article.CreatedBy),
		sql.Named("ModifiedBy", article.ModifiedBy))

	if err != nil {
		return err
	}
	err = row.Scan(&article)
	if err != nil {
		return err
	}

	return nil
}
