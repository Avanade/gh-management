package contributionarea

import (
	"database/sql"
	"main/model"
	"main/repository"
	"time"
)

type contributionAreaRepository struct {
	repository.Database
}

// GetAll implements ContributionAreaRepository.
func (r *contributionAreaRepository) GetAll() ([]model.ContributionArea, error) {
	var contributionAreas []model.ContributionArea
	rows, err := r.Query("[dbo].[usp_ContributionArea_Select]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var contributionArea model.ContributionArea

		contributionArea.ID = v["Id"].(int64)
		contributionArea.Name = v["Name"].(string)
		contributionArea.Created = v["Created"].(time.Time)
		contributionArea.CreatedBy = v["CreatedBy"].(string)
		contributionArea.Modified = v["Modified"].(time.Time)
		contributionArea.ModifiedBy = v["ModifiedBy"].(string)

		contributionAreas = append(contributionAreas, contributionArea)
	}

	return contributionAreas, nil
}

// GetByOption implements ContributionAreaRepository.
func (r *contributionAreaRepository) GetByOption(offset int, filter int, orderby string, ordertype string, search string) ([]model.ContributionArea, error) {
	var contributionAreas []model.ContributionArea
	rows, err := r.Query("[dbo].[usp_ContributionArea_Select_ByOption]",
		sql.Named("Offset", offset),
		sql.Named("Filter", filter),
		sql.Named("OrderBy", orderby),
		sql.Named("OrderType", ordertype),
		sql.Named("Search", search))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var contributionArea model.ContributionArea

		contributionArea.ID = v["Id"].(int64)
		contributionArea.Name = v["Name"].(string)
		contributionArea.Created = v["Created"].(time.Time)
		contributionArea.CreatedBy = v["CreatedBy"].(string)
		contributionArea.Modified = v["Modified"].(time.Time)
		if v["ModifiedBy"] != nil {
			contributionArea.ModifiedBy = v["ModifiedBy"].(string)
		}
		contributionAreas = append(contributionAreas, contributionArea)
	}

	return contributionAreas, nil
}

// GetByID implements ContributionAreaRepository.
func (r *contributionAreaRepository) GetByID(id int64) (*model.ContributionArea, error) {
	var contributionArea model.ContributionArea
	row, err := r.QueryRow("[dbo].[usp_ContributionArea_Select_ById]",
		sql.Named("Id", id))
	if err != nil {
		return nil, err
	}

	var modifiedBy sql.NullString

	err = row.Scan(
		&contributionArea.ID,
		&contributionArea.Name,
		&contributionArea.Created,
		&contributionArea.CreatedBy,
		&contributionArea.Modified,
		&modifiedBy,
	)
	if err != nil {
		return nil, err
	}

	if modifiedBy.Valid {
		contributionArea.ModifiedBy = modifiedBy.String
	}

	return &contributionArea, nil
}

// Create implements ContributionAreaRepository.
func (r *contributionAreaRepository) Create(contributionArea *model.ContributionArea) (*model.ContributionArea, error) {
	result, err := r.QueryRow("[dbo].[usp_ContributionArea_Insert]",
		sql.Named("Name", contributionArea.Name),
		sql.Named("CreatedBy", contributionArea.CreatedBy))
	if err != nil {
		return nil, err
	}
	err = result.Scan(
		&contributionArea.ID,
		&contributionArea.Created)
	if err != nil {
		return nil, err
	}
	return contributionArea, nil
}

// Update implements ContributionAreaRepository.
func (r *contributionAreaRepository) Update(id int64, contributionArea *model.ContributionArea) (*model.ContributionArea, error) {
	row, err := r.QueryRow("[dbo].[usp_ContributionArea_Update]",
		sql.Named("Id", id),
		sql.Named("Name", contributionArea.Name),
		sql.Named("ModifiedBy", contributionArea.ModifiedBy))
	if err != nil {
		return nil, err
	}
	err = row.Scan(
		&contributionArea.Modified)
	if err != nil {
		return nil, err
	}
	return contributionArea, nil
}

// GetTotal implements ContributionAreaRepository.
func (r *contributionAreaRepository) GetTotal() (int, error) {
	var total int
	row, err := r.QueryRow("[dbo].[usp_ContributionArea_TotalCount]")
	if err != nil {
		return 0, err
	}
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func NewContributionAreaRepository(database repository.Database) ContributionAreaRepository {
	return &contributionAreaRepository{database}
}
