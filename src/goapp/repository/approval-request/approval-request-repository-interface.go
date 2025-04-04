package approvalRequest

type ApprovalRequestRepository interface {
	Insert(approver, description, requestor string) (id int64, err error)
	UpdateApprovalSystemGUID(requestId int64, approvalSystemGUID string) error
}
