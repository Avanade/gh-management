package adoOrganization

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
)

type adoOrganizationRepository struct {
	*db.Database
}

func NewAdoOrganizationRepository(db *db.Database) AdoOrganizationRepository {
	return &adoOrganizationRepository{db}
}

func (r *adoOrganizationRepository) Insert(adoOrganization *model.AdoOrganizationRequest) (int, error) {
	row, err := r.QueryRow("usp_AdoOrganization_Insert",
		sql.Named("Name", adoOrganization.Name),
		sql.Named("Purpose", adoOrganization.Purpose),
		sql.Named("CreatedBy", adoOrganization.CreatedBy),
	)
	if err != nil {
		return 0, err
	}

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *adoOrganizationRepository) SelectByUser(user string) ([]model.AdoOrganizationRequest, error) {
	rows, err := r.Query("usp_AdoOrganization_Select_ByUsername", sql.Named("Username", user))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	adoOrganizations, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	var result []model.AdoOrganizationRequest
	for _, adoOrganization := range adoOrganizations {
		result = append(result, model.AdoOrganizationRequest{
			Id:        int(adoOrganization["Id"].(int64)),
			Name:      adoOrganization["Name"].(string),
			Purpose:   adoOrganization["Purpose"].(string),
			CreatedBy: adoOrganization["CreatedBy"].(string),
			Created:   adoOrganization["Created"].(time.Time),
		})
	}

	return result, nil
}
