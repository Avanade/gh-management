package approvalRequest

import (
	"database/sql"
	db "main/infrastructure/database"
)

type approvalRequestRepository struct {
	*db.Database
}

func NewApprovalRequestRepository(db *db.Database) ApprovalRequestRepository {
	return &approvalRequestRepository{db}
}

func (r *approvalRequestRepository) Insert(approver, description, requestor string) (id int64, err error) {
	row, err := r.Query("usp_ApprovalRequest_Select_Insert",
		sql.Named("ApproverUserPrincipalName", approver),
		sql.Named("Name", description),
		sql.Named("CreatedBy", requestor))
	if err != nil {
		return 0, err
	}
	defer row.Close()

	row.Next()
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *approvalRequestRepository) UpdateApprovalSystemGUID(requestId int64, approvalSystemGUID string) error {
	_, err := r.Query("usp_ApprovalRequest_Update_ApprovalSystemGUID",
		sql.Named("Id", requestId),
		sql.Named("ApprovalSystemGUID", approvalSystemGUID))
	return err
}
