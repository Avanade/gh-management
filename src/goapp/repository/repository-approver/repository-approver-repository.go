package repositoryApprover

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
)

type repositoryApproverRepository struct {
	*db.Database
}

func NewRepostioryApproverRepository(db *db.Database) RepositoryApproverRepository {
	return &repositoryApproverRepository{db}
}

func (r *repositoryApproverRepository) Insert(approver *model.RepositoryApprover) error {
	err := r.Execute("usp_RepositoryApprover_Insert",
		sql.Named("RepositoryApprovalTypeId", approver.ApprovalTypeId),
		sql.Named("ApproverUserPrincipalName", approver.ApproverEmail))
	return err
}

func (r *repositoryApproverRepository) SelectByApprovalTypeId(approvalTypeId int) ([]model.RepositoryApprover, error) {
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
