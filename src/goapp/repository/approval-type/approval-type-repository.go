package approvalType

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type approvalTypeRepository struct {
	*db.Database
}

func NewApprovalTypeRepository(db *db.Database) ApprovalTypeRepository {
	return &approvalTypeRepository{db}
}

func (r *approvalTypeRepository) Select() ([]model.ApprovalType, error) {
	rows, err := r.Query("usp_RepositoryApprovalType_Select")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	approvalTypes, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	var result []model.ApprovalType
	for _, approvalType := range approvalTypes {
		result = append(result, model.ApprovalType{
			Id:         int(approvalType["Id"].(int64)),
			Name:       approvalType["Name"].(string),
			IsActive:   approvalType["IsActive"].(bool),
			IsArchived: approvalType["IsArchived"].(bool),
		})
	}

	return result, nil
}

func (r *approvalTypeRepository) SelectByOption(opt model.FilterOptions) ([]model.ApprovalType, error) {
	rows, err := r.Query("usp_RepositoryApprovalType_Select_ByOption",
		sql.Named("Offset", opt.Offset),
		sql.Named("Filter", opt.Filter),
		sql.Named("OrderBy", opt.Orderby),
		sql.Named("OrderType", opt.Ordertype),
		sql.Named("Search", opt.Search),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	approvalTypes, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	var result []model.ApprovalType
	for _, approvalType := range approvalTypes {
		result = append(result, model.ApprovalType{
			Id:         int(approvalType["Id"].(int64)),
			Name:       approvalType["Name"].(string),
			IsActive:   approvalType["IsActive"].(bool),
			IsArchived: approvalType["IsArchived"].(bool),
		})
	}

	return result, nil
}

func (r *approvalTypeRepository) Total() (int64, error) {
	row, err := r.QueryRow("usp_RepositoryApprovalType_TotalCount")
	if err != nil {
		return 0, err
	}

	var total int64
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
