package category

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
)

type categoryRepository struct {
	*db.Database
}

func NewCategoryRepository(db *db.Database) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) Insert(category *model.Category) (int64, error) {
	row, err := r.QueryRow("usp_GuidanceCategory_Insert",
		sql.Named("Name", category.Name),
		sql.Named("CreatedBy", category.CreatedBy),
		sql.Named("ModifiedBy", category.ModifiedBy),
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

func (r *categoryRepository) Select() ([]model.Category, error) {
	var categories []model.Category
	rows, err := r.Query("usp_GuidanceCategory_Select")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var category model.Category

		category.Id = v["Id"].(int64)
		category.Name = v["Name"].(string)
		category.Created = v["Created"].(time.Time)
		category.CreatedBy = v["CreatedBy"].(string)
		if v["Modified"] != nil {
			category.Modified = v["Modified"].(time.Time)
		}
		if v["ModifiedBy"] != nil {
			category.ModifiedBy = v["ModifiedBy"].(string)
		}

		categories = append(categories, category)
	}
	return categories, nil
}

func (r *categoryRepository) SelectById(id int64) (*model.Category, error) {
	var category model.Category

	row, err := r.QueryRow("usp_GuidanceCategory_Select_ById", sql.Named("Id", id))
	if err != nil {
		return nil, err
	}

	var modified sql.NullTime
	var modifiedBy sql.NullString

	err = row.Scan(
		&category.Id,
		&category.Name,
		&category.Created,
		&category.CreatedBy,
		&modified,
		&modifiedBy,
	)
	if err != nil {
		return nil, err
	}

	if modified.Valid {
		category.Modified = modified.Time
	}
	if modifiedBy.Valid {
		category.ModifiedBy = modifiedBy.String
	}

	return &category, nil
}

func (r *categoryRepository) Update(category *model.Category) error {
	_, err := r.QueryRow("[dbo].[usp_GuidanceCategory_Update]",
		sql.Named("Id", category.Id),
		sql.Named("Name", category.Name),
		sql.Named("CreatedBy", category.CreatedBy),
		sql.Named("ModifiedBy", category.ModifiedBy))

	if err != nil {
		return err
	}

	return nil
}
