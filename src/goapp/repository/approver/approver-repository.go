package approver

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type approverRepository struct {
	*db.Database
}

func NewApproverRepository(db *db.Database) ApproverRepository {
	return &approverRepository{db}
}

func (r *approverRepository) SelectByApprovalTypeId(approvalTypeId int) ([]model.RepositoryApprover, error) {
	rows, err := r.Query("usp_RepositoryApprover_Select_ByApprovalTypeId", sql.Named("RepositoryApprovalTypeId", approvalTypeId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	approvers, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	var result []model.RepositoryApprover
	for _, approver := range approvers {
		result = append(result, model.RepositoryApprover{
			ApprovalTypeId: int(approver["RepositoryApprovalTypeId"].(int64)),
			ApproverEmail:  approver["ApproverUserPrincipalName"].(string),
			ApproverName:   approver["ApproverName"].(string),
		})
	}

	return result, nil
}
