package adoOrganizationApprovalRequest

import (
	"database/sql"
	db "main/infrastructure/database"
)

type adoOrganizationApprovalRequestRepository struct {
	*db.Database
}

func NewAdoOrganizationApprovalRequestRepository(db *db.Database) AdoOrganizationApprovalRequestRepository {
	return &adoOrganizationApprovalRequestRepository{db}
}

func (r *adoOrganizationApprovalRequestRepository) Insert(adoOrganizationId int, approvalRequestId int64) error {
	_, err := r.Query("usp_AdoOrganizationApprovalRequest_Insert",
		sql.Named("AdoOrganizationId", adoOrganizationId),
		sql.Named("ApprovalRequestId", approvalRequestId))
	if err != nil {
		return err
	}

	return nil
}
