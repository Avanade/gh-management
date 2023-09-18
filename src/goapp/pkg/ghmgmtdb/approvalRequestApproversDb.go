package ghmgmt

type ApprovalRequestApprover struct {
	ApprovalRequestId int
	ApproverEmail     string
}

func InsertApprovalRequestApprover(approvalRequestApprover ApprovalRequestApprover) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalRequestId": approvalRequestApprover.ApprovalRequestId,
		"ApproverEmail":     approvalRequestApprover.ApproverEmail,
	}

	_, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalRequestApprovers_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func GetApprovalRequestApproversByApprovalRequestId(approvalRequestId int) ([]ApprovalRequestApprover, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalRequestId": approvalRequestId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalRequestApprovers_Select_ByApprovalRequestId", param)
	if err != nil {
		return nil, err
	}

	var approvalRequestApprovers []ApprovalRequestApprover

	for _, v := range result {
		approvalRequestApprover := ApprovalRequestApprover{
			ApprovalRequestId: int(v["ApprovalTypeId"].(int64)),
			ApproverEmail:     v["ApproverEmail"].(string),
		}

		approvalRequestApprovers = append(approvalRequestApprovers, approvalRequestApprover)
	}

	return approvalRequestApprovers, nil
}
