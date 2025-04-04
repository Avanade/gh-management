package communityApprover

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type communityApproverRepository struct {
	*db.Database
}

func NewCommunityApproverRepository(db *db.Database) CommunityApproverRepository {
	return &communityApproverRepository{db}
}

func (r *communityApproverRepository) GetByCategory(category string) ([]model.CommunityApprover, error) {
	row, err := r.Query("usp_CommunityApprover_Select_Active",
		sql.Named("GuidanceCategory", category))
	if err != nil {
		return nil, err
	}
	defer row.Close()

	approvers, err := r.RowsToMap(row)
	if err != nil {
		return nil, err
	}

	var communityApprovers []model.CommunityApprover
	for _, approver := range approvers {
		communityApprovers = append(communityApprovers, model.CommunityApprover{
			Id:                        approver["Id"].(int64),
			ApproverUserPrincipalName: approver["ApproverUserPrincipalName"].(string),
			GuidanceCategory:          category,
		})
	}

	return communityApprovers, nil
}
