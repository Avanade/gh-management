package adoOrganizationApprovalRequest

type AdoOrganizationApprovalRequestRepository interface {
	Insert(adoOrganizationId int, approvalRequestId int64) error
}
