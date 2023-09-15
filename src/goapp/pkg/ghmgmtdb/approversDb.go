package ghmgmt

type Approver struct {
	ApprovalTypeId int
	ApproverEmail  string
}

func InsertApprover(approver Approver) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalTypeId": approver.ApprovalTypeId,
		"ApproverEmail":  approver.ApproverEmail,
	}

	_, err := db.ExecuteStoredProcedureWithResult("PR_Approvers_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func DeleteApproverByApprovalTypeId(approvalTypeId int) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalTypeId": approvalTypeId,
	}

	_, err := db.ExecuteStoredProcedureWithResult("PR_Approvers_Delete_ByApprovalTypeId", param)
	if err != nil {
		return err
	}
	return nil
}

func GetApproversByApprovalTypeId(approvalTypeId int) ([]Approver, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalTypeId": approvalTypeId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_Approvers_Select_ByApprovalTypeId", param)
	if err != nil {
		return nil, err
	}

	var approvers []Approver

	for _, v := range result {
		approver := Approver{
			ApprovalTypeId: int(v["ApprovalTypeId"].(int64)),
			ApproverEmail:  v["ApproverEmail"].(string),
		}

		approvers = append(approvers, approver)
	}

	return approvers, nil
}
