package ghmgmt

type GitHubCopilotDto struct {
	Region                    int    `json:"region"`
	GitHubId                  int64  `json:"githubId"`
	GitHubUsername            string `json:"githubUsername"`
	Username                  string
	Id                        int64
	ApproverUserPrincipalName []string
	RegionName                string `json:"regionName"`
	RequestId                 []int64
}

func GitHubCopilotInsert(body GitHubCopilotDto) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"Region":         body.Region,
		"GitHubId":       body.GitHubId,
		"GitHubUsername": body.GitHubUsername,
		"Username":       body.Username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_GitHubCopilot_Insert", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GitHubCopilotApprovalInsert(githubCopilotId int, requestId int64) error {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{

		"GitHubCopilotId": githubCopilotId,
		"RequestId":       requestId,
	}

	_, err := db.ExecuteStoredProcedure("dbo.PR_GitHubCopilotApprovalRequests_Insert", params)
	if err != nil {
		return err
	}

	return nil
}

func GetAllGitHubCopilotRequest(username string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Username": username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_GitHubCopilot_SelectByUser", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

func GetGitHubCopilotApprovalRequest(id int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_GitHubCopilotApprovalRequests_SelectByCopilotId", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

func GetGitHubCopilotbyGUID(id string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalSystemGUID": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_GitHubCopilot_SelectByGUID", param)
	if err != nil {
		return nil, err
	}

	return result, err
}
