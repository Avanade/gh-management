package osscontributionsponsor

import (
	"database/sql"
	"main/model"
	"main/repository"
)

type ossContributionSponsorRepository struct {
	repository.Database
}

// Create implements OSSContributionSponsorRepository.
func (r *ossContributionSponsorRepository) Create(ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error) {
	result, err := r.QueryRow("[dbo].[usp_OSSContributionSponsor_Insert]",
		sql.Named("Name", ossContributionSponsor.Name),
		sql.Named("IsArchived", ossContributionSponsor.IsArchived))
	if err != nil {
		return nil, err
	}
	err = result.Scan(&ossContributionSponsor.ID)
	if err != nil {
		return nil, err
	}
	return ossContributionSponsor, nil
}

// GetAll implements OSSContributionSponsorRepository.
func (r *ossContributionSponsorRepository) GetAll() ([]model.OSSContributionSponsor, error) {
	var ossContributionSponsors []model.OSSContributionSponsor
	rows, err := r.Query("[dbo].[usp_OSSContributionSponsor_Select]")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var ossContributionSponsor model.OSSContributionSponsor

		ossContributionSponsor.ID = v["Id"].(int64)
		ossContributionSponsor.Name = v["Name"].(string)
		ossContributionSponsor.IsArchived = v["IsArchived"].(bool)

		ossContributionSponsors = append(ossContributionSponsors, ossContributionSponsor)
	}

	return ossContributionSponsors, nil
}

// GetByID implements OSSContributionSponsorRepository.
func (r *ossContributionSponsorRepository) GetByIsArchived(isArchived bool) ([]model.OSSContributionSponsor, error) {
	var ossContributionSponsors []model.OSSContributionSponsor
	rows, err := r.Query("[dbo].[usp_OSSContributionSponsor_Select_ByIsArchived]",
		sql.Named("IsArchived", isArchived))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapRows, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	for _, v := range mapRows {
		var ossContributionSponsor model.OSSContributionSponsor
		ossContributionSponsor.ID = v["Id"].(int64)
		ossContributionSponsor.Name = v["Name"].(string)
		ossContributionSponsor.IsArchived = v["IsArchived"].(bool)
		ossContributionSponsors = append(ossContributionSponsors, ossContributionSponsor)
	}

	return ossContributionSponsors, nil
}

// Update implements OSSContributionSponsorRepository.
func (r *ossContributionSponsorRepository) Update(id int64, ossContributionSponsor *model.OSSContributionSponsor) (*model.OSSContributionSponsor, error) {
	err := r.Execute("[dbo].[usp_OSSContributionSponsor_Update]",
		sql.Named("Id", id),
		sql.Named("Name", ossContributionSponsor.Name),
		sql.Named("IsArchived", ossContributionSponsor.IsArchived))
	if err != nil {
		return nil, err
	}
	return ossContributionSponsor, nil
}

func NewOSSContributionSponsorRepository(db repository.Database) OSSContributionSponsorRepository {
	return &ossContributionSponsorRepository{
		Database: db,
	}
}
