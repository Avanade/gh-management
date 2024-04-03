package ghmgmt

import (
	"strconv"
	"strings"
)

type Organization struct {
	Id          int64
	RegionId    int `json:"region"`
	RegionName  string
	ClientName  string `json:"clientName"`
	ProjectName string `json:"projectName"`
	WBS         string `json:"wbs"`
	Username    string
	Approvers   []string
	RequestIds  []int64
}

func OrganizationInsert(body Organization) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"Region":      body.RegionId,
		"ClientName":  body.ClientName,
		"ProjectName": body.ProjectName,
		"WBS":         body.WBS,
		"CreatedBy":   body.Username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Organizations_Insert", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetRegionalOrganizationById(id int) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RegionalOrganizations_SelectById", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetAllRegionalOrganizations() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_RegionalOrganizations_SelectAll", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func OrganizationUpdateApprovalStatus(id int64, status int) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":               id,
		"ApprovalStatusID": status,
	}
	db.ExecuteStoredProcedure("PR_Organizations_UpdateStatus", param)
}

func GetAllOrganizationRequest(username string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Username": username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Organizations_SelectByUser", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

func GetOrganizationApprovalRequest(id int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_OrganizationsApprovalRequests_SelectByOrgId", param)
	if err != nil {
		return nil, err
	}

	return result, err
}

func RegionalOrganizationInsert(id int64, name string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":   id,
		"Name": name,
	}

	_, err := db.ExecuteStoredProcedure("dbo.PR_RegionalOrganizations_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

func OrganizationApprovalInsert(organizationId int, requestId int64) error {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{

		"OrganizationId": organizationId,
		"RequestId":      requestId,
	}

	_, err := db.ExecuteStoredProcedure("dbo.PR_OrganizationsApprovalRequests_Insert", params)
	if err != nil {
		return err
	}

	return nil
}

func GetFailedCommunityApprovalRequestNewOrganizations() []Organization {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_CommunityApprovals_Select_FailedRequestOrganizations", nil)

	var organizations []Organization

	for _, v := range result {
		data := Organization{
			Id:          v["Id"].(int64),
			RegionId:    int(v["RegionId"].(int64)),
			RegionName:  v["RegionName"].(string),
			ClientName:  v["ClientName"].(string),
			ProjectName: v["ProjectName"].(string),
			WBS:         v["WBS"].(string),
			Username:    v["Username"].(string),
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

		organizations = append(organizations, data)
	}

	return organizations
}
