package adoOrganizationApprovalRequest

type AdoOrganizationApprovalRequestService interface {
	Insert(adoOrganizationId int, approvalRequestId int64) error
}
