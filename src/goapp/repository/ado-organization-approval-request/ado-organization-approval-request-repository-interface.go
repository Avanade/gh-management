package adoOrganizationApprovalRequest

import "main/model"

type AdoOrganizationApprovalRequestRepository interface {
	Insert(adoOrganizationId int, approvalRequestId int64) error
	SelectByAdoOrganizationId(adoOrganizationId int64) ([]model.ApprovalRequest, error)
}
