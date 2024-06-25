package ghmgmt

type Approver struct {
	ApprovalTypeId int
	ApproverEmail  string
	ApproverName   string
}

func InsertApprover(approver Approver) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryApprovalTypeId":  approver.ApprovalTypeId,
		"ApproverUserPrincipalName": approver.ApproverEmail,
	}

	_, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprover_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func DeleteApproverByApprovalTypeId(approvalTypeId int) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryApprovalTypeId": approvalTypeId,
	}

	_, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprover_Delete", param)
	if err != nil {
		return err
	}
	return nil
}

func GetApproversByApprovalTypeId(approvalTypeId int) ([]Approver, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryApprovalTypeId": approvalTypeId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApprover_Select_ByApprovalTypeId", param)
	if err != nil {
		return nil, err
	}

	var approvers []Approver

	for _, v := range result {
		approver := Approver{
			ApprovalTypeId: int(v["RepositoryApprovalTypeId"].(int64)),
			ApproverEmail:  v["ApproverUserPrincipalName"].(string),
			ApproverName:   v["ApproverName"].(string),
		}

		approvers = append(approvers, approver)
	}

	return approvers, nil
}
