package ghmgmt

import (
	"strconv"
	"strings"
)

type GitHubCopilot struct {
	Region         int    `json:"region"`
	GitHubId       int64  `json:"githubId"`
	GitHubUsername string `json:"githubUsername"`
	Username       string
	Id             int64
	RegionName     string `json:"regionName"`
	Approvers      []string
	RequestIds     []int64
}

func GitHubCopilotInsert(body GitHubCopilot) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"RegionalOrganizationId": body.Region,
		"GitHubId":               body.GitHubId,
		"GitHubUsername":         body.GitHubUsername,
		"Username":               body.Username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_GitHubCopilot_Insert", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("usp_GitHubCopilot_Select_ByUsername", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("usp_GitHubCopilot_Select_ByApprovalSystemGUID", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

func GitHubCopilotGetPendingByUserAndOrganization(body GitHubCopilot) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Username":     body.Username,
		"Organization": body.RegionName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_GitHubCopilot_Select_Pending", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetFailedCommunityApprovalRequestGitHubCoPilots() []GitHubCopilot {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Select_FailedRequestGitHubCopilot", nil)

	var gitHubCopilots []GitHubCopilot

	for _, v := range result {
		githubId, _ := strconv.ParseInt(v["GitHubId"].(string), 0, 64)
		data := GitHubCopilot{
			Region:         int(v["RegionId"].(int64)),
			GitHubId:       githubId,
			GitHubUsername: v["GitHubUsername"].(string),
			Username:       v["Username"].(string),
			Id:             v["Id"].(int64),
			RegionName:     v["RegionName"].(string),
			Approvers:      []string{},
			RequestIds:     []int64{},
		}

		if v["Approvers"] != nil {
			approversStr := v["Approvers"].(string)
			data.Approvers = strings.Split(approversStr, ",")
		}

		if v["RequestIds"] != nil {
			requestIdsStr := v["RequestIds"].(string)
			requestIds := strings.Split(requestIdsStr, ",")

			for _, v := range requestIds {
				requestId, err := strconv.ParseInt(v, 0, 64)
				if err != nil {
					continue
				}
				data.RequestIds = append(data.RequestIds, int64(requestId))
			}
		}

		gitHubCopilots = append(gitHubCopilots, data)
	}

	return gitHubCopilots
}
