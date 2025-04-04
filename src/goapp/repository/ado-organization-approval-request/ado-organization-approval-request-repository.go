package adoOrganizationApprovalRequest

import (
	"database/sql"
	db "main/infrastructure/database"
	"main/model"
	"time"
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

func (r *adoOrganizationApprovalRequestRepository) SelectByAdoOrganizationId(adoOrganizationId int64) ([]model.ApprovalRequest, error) {
	rows, err := r.Query("usp_AdoOrganizationApprovalRequest_Select_ByAdoOrganizationId",
		sql.Named("AdoOrganizationId", adoOrganizationId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	approvalRequests, err := r.RowsToMap(rows)
	if err != nil {
		return nil, err
	}

	var result []model.ApprovalRequest
	for _, approvalRequest := range approvalRequests {
		var req = model.ApprovalRequest{
			Id:                  approvalRequest["Id"].(int64),
			ApprovalStatus:      approvalRequest["ApprovalStatus"].(string),
			ApprovalDescription: approvalRequest["ApprovalDescription"].(string),
		}

		if approvalRequest["ApprovalDate"] != nil {
			req.ApprovalDate = approvalRequest["ApprovalDate"].(time.Time)
		}

		if approvalRequest["ApproverUserPrincipalName"] != nil {
			req.ApproverPrincipalName = approvalRequest["ApproverUserPrincipalName"].(string)
		}

		if approvalRequest["ApprovalRemarks"] != nil {
			req.ApprovalRemarks = approvalRequest["ApprovalRemarks"].(string)
		}

		result = append(result, req)
	}

	return result, nil
}
