package ghmgmt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type InsertGHCPPremiumBudgetAllocationInput struct {
	OrganizationGitHubId    int64
	OrganizationGitHubLogin string
	UserGitHubId            int64
	UserGitHubLogin         string
	Amount                  int64
	CreatedBy               string
}

type InsertGHCPPremiumBudgetAllocationApprovalRequestInput struct {
	GHCPPremiumBudgetAllocationID int64
	ApprovalRequestId             int64
}

func InsertGHCPPremiumBudgetAllocation(input InsertGHCPPremiumBudgetAllocationInput) (id int64, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"OrganizationGitHubID":    input.OrganizationGitHubId,
		"OrganizationGitHubLogin": input.OrganizationGitHubLogin,
		"UserGitHubID":            input.UserGitHubId,
		"UserGitHubLogin":         input.UserGitHubLogin,
		"Amount":                  input.Amount,
		"CreatedBy":               input.CreatedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_GHCPPremiumBudgetAllocation_Insert", param)
	if err != nil {
		return 0, err
	}

	id, err = strconv.ParseInt(fmt.Sprint(result[0]["Id"]), 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func InsertGHCPPremiumBudgetAllocationApprovalRequest(input InsertGHCPPremiumBudgetAllocationApprovalRequestInput) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GHCPPremiumBudgetAllocationId": input.GHCPPremiumBudgetAllocationID,
		"ApprovalRequestId":             input.ApprovalRequestId,
	}

	_, err := db.ExecuteStoredProcedureWithResult("usp_GHCPPremiumBudgetAllocationApprovalRequest_Insert", param)
	if err != nil {
		return err
	}

	return nil
}

type GHCPPremiumBudgetAllocation struct {
	Id                      int64
	OrganizationGitHubID    int64
	OrganizationGitHubLogin string
	UserGitHubID            int64
	UserGitHubLogin         string
	Amount                  int64
	Created                 time.Time
	CreatedBy               string
}

func GetGHCPPremiumBudgetAllocationRequestByUsername(username string) (ghcpPremiuBudgetAllocationRequests []GHCPPremiumBudgetAllocation, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Username": username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_GHCPPremiumBudgetAllocation_Select_ByUsername", param)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		ghcpPremiuBudgetAllocationRequests = append(ghcpPremiuBudgetAllocationRequests, GHCPPremiumBudgetAllocation{
			Id:                      v["Id"].(int64),
			OrganizationGitHubID:    v["OrganizationGitHubID"].(int64),
			OrganizationGitHubLogin: v["OrganizationGitHubLogin"].(string),
			UserGitHubID:            v["UserGitHubID"].(int64),
			UserGitHubLogin:         v["UserGitHubLogin"].(string),
			Amount:                  v["Amount"].(int64),
			Created:                 v["Created"].(time.Time),
			CreatedBy:               v["CreatedBy"].(string),
		})
	}

	return ghcpPremiuBudgetAllocationRequests, nil
}

type GHCPPremiumBudgetAllocationApprovalRequest struct {
	Id           int64
	Approver     string
	Status       string
	Remarks      *string
	ResponseDate *time.Time
	Description  *string
}

func GetGHCPPremiumBudgetAllocationApprovalRequest(id int64) (ghcpPremiumBudgetAllocationApprovalRequest []GHCPPremiumBudgetAllocationApprovalRequest, err error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GHCPPremiumBudgetAllocationId": id,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_GHCPPremiumBudgetAllocationApprovalRequest_Select_ByGHCPPremiumBudgetAllocationId", param)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		var remarks *string
		if v["ApprovalRemarks"] != nil {
			remarksValue := v["ApprovalRemarks"].(string)
			remarks = &remarksValue
		}

		var responseDate *time.Time
		if v["ApprovalDate"] != nil {
			responseDateValue := v["ApprovalDate"].(time.Time)
			responseDate = &responseDateValue
		}

		var description *string
		if v["ApprovalDescription"] != nil {
			descriptionValue := v["ApprovalDescription"].(string)
			description = &descriptionValue
		}

		ghcpPremiumBudgetAllocationApprovalRequest = append(ghcpPremiumBudgetAllocationApprovalRequest, GHCPPremiumBudgetAllocationApprovalRequest{
			Id:           v["Id"].(int64),
			Approver:     v["ApproverUserPrincipalName"].(string),
			Status:       v["ApprovalStatus"].(string),
			Remarks:      remarks,
			ResponseDate: responseDate,
			Description:  description,
		})
	}

	return ghcpPremiumBudgetAllocationApprovalRequest, nil
}

type FailedCommunityApprovalRequestGHCPPremiumBudgetAllocation struct {
	GHCPPremiumBudgetAllocation
	Approvers  []string
	RequestIds []int64
}

func GetFailedCommunityApprovalRequestGHCPPremiumBudgetAllocation() (failedCommunityApprovalRequestGHCPPremiumBudgetAllocations []FailedCommunityApprovalRequestGHCPPremiumBudgetAllocation) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Select_FailedRequestGHCPPremiumBudgetAllocation", nil)

	for _, v := range result {
		failedCommunityApprovalRequestGHCPPremiumBudgetAllocation := FailedCommunityApprovalRequestGHCPPremiumBudgetAllocation{
			GHCPPremiumBudgetAllocation: GHCPPremiumBudgetAllocation{
				Id:                      v["Id"].(int64),
				OrganizationGitHubID:    v["OrganizationGitHubId"].(int64),
				OrganizationGitHubLogin: v["OrganizationGitHubLogin"].(string),
				UserGitHubID:            v["UserGitHubId"].(int64),
				UserGitHubLogin:         v["UserGitHubLogin"].(string),
				Amount:                  v["Amount"].(int64),
				Created:                 v["Requested"].(time.Time),
				CreatedBy:               v["RequestedBy"].(string),
			},
		}
		if v["Approvers"] != nil {
			approversStr := v["Approvers"].(string)
			failedCommunityApprovalRequestGHCPPremiumBudgetAllocation.Approvers = strings.Split(approversStr, ",")
		}

		if v["RequestIds"] != nil {
			requestIdsStr := v["RequestIds"].(string)
			requestIds := strings.Split(requestIdsStr, ",")

			for _, v := range requestIds {
				requestId, err := strconv.ParseInt(v, 0, 64)
				if err != nil {
					continue
				}
				failedCommunityApprovalRequestGHCPPremiumBudgetAllocation.RequestIds = append(failedCommunityApprovalRequestGHCPPremiumBudgetAllocation.RequestIds, int64(requestId))
			}
		}

		failedCommunityApprovalRequestGHCPPremiumBudgetAllocations = append(failedCommunityApprovalRequestGHCPPremiumBudgetAllocations, failedCommunityApprovalRequestGHCPPremiumBudgetAllocation)
	}

	return failedCommunityApprovalRequestGHCPPremiumBudgetAllocations
}
