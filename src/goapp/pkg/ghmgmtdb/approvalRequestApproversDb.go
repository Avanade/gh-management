package ghmgmt

import "time"

type ApprovalRequestApprover struct {
	ApprovalRequestId int
	ApproverEmail     string
}

type ProjectApprovalApprovers struct {
	Id                         int64
	ProjectId                  int64
	ProjectName                string
	ProjectDescription         string
	RequesterName              string
	RequesterGivenName         string
	RequesterSurName           string
	RequesterUserPrincipalName string
	ApprovalTypeId             int
	ApprovalType               string
	Approvers                  []string
	ApprovalDescription        string
	RequestStatus              string
	ApprovalDate               time.Time
	ApprovalRemarks            string
	ConfirmAvaIP               bool
	ConfirmEnabledSecurity     bool
	ConfirmNotClientProject    bool
	NewContribution            string
	OSSsponsor                 string
	OfferingsAssets            string
	WillBeCommercialVersion    string
	OSSContributionInformation string
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
			ApprovalRequestId: int(v["ApprovalRequestId"].(int64)),
			ApproverEmail:     v["ApproverEmail"].(string),
		}

		approvalRequestApprovers = append(approvalRequestApprovers, approvalRequestApprover)
	}

	return approvalRequestApprovers, nil
}

func PopulateApprovalRequestApproversProjectApprovalsByProject(projectId int64, requestedBy string) ([]ProjectApprovalApprovers, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId":   projectId,
		"RequestedBy": requestedBy,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_ApprovalRequestApprovers_Populate", param)
	if err != nil {
		return nil, err
	}

	var projectApprovalApprovers []ProjectApprovalApprovers

	for _, v := range result {
		projectApprovalApprover := ProjectApprovalApprovers{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApprovalTypeId:             int(v["ApprovalTypeId"].(int64)),
			ApprovalType:               v["ApprovalType"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			RequestStatus:              v["RequestStatus"].(string),
			ConfirmAvaIP:               v["ConfirmAvaIP"].(bool),
			ConfirmEnabledSecurity:     v["ConfirmEnabledSecurity"].(bool),
			ConfirmNotClientProject:    v["ConfirmNotClientProject"].(bool),
		}

		if v["ApprovalDate"] != nil {
			projectApprovalApprover.ApprovalDate = v["ApprovalDate"].(time.Time)
		}

		if v["ApprovalRemarks"] != nil {
			projectApprovalApprover.ApprovalRemarks = v["ApprovalRemarks"].(string)
		}

		if v["newcontribution"] != nil {
			projectApprovalApprover.NewContribution = v["newcontribution"].(string)
		}

		if v["Willbecommercialversion"] != nil {
			projectApprovalApprover.WillBeCommercialVersion = v["Willbecommercialversion"].(string)
		}

		if v["Avanadeofferingsassets"] != nil {
			projectApprovalApprover.OfferingsAssets = v["Avanadeofferingsassets"].(string)
		}

		if v["OSSsponsor"] != nil {
			projectApprovalApprover.OSSsponsor = v["OSSsponsor"].(string)
		}

		if v["OSSContributionInformation"] != nil {
			projectApprovalApprover.OSSContributionInformation = v["OSSContributionInformation"].(string)
		}

		projectApprovalApprovers = append(projectApprovalApprovers, projectApprovalApprover)
	}

	return projectApprovalApprovers, nil
}

func RequestProjectApprovals(projectId int64, requestedBy string) ([]ProjectApprovalApprovers, error) {
	projectApprovalApprovers, err := PopulateApprovalRequestApproversProjectApprovalsByProject(projectId, requestedBy)
	if err != nil {
		return nil, err
	}

	for index, projectApprovalApprover := range projectApprovalApprovers {
		approvalRequestApprovers, err := GetApprovalRequestApproversByApprovalRequestId(int(projectApprovalApprover.Id))
		if err != nil {
			return nil, err
		}

		for _, approvalRequestApprover := range approvalRequestApprovers {
			projectApprovalApprovers[index].Approvers = append(projectApprovalApprovers[index].Approvers, approvalRequestApprover.ApproverEmail)
		}
	}

	return projectApprovalApprovers, nil
}

func ReprocessFailedProjectApprovals() ([]ProjectApprovalApprovers, error) {
	projectApprovals := GetFailedProjectApprovalRequests()

	var projectApprovalApprovers []ProjectApprovalApprovers

	for _, v := range projectApprovals {
		projectApprovalApprover := ProjectApprovalApprovers{
			Id:                         v.Id,
			ProjectId:                  v.ProjectId,
			ProjectName:                v.ProjectName,
			ProjectDescription:         v.ProjectDescription,
			RequesterName:              v.RequesterName,
			RequesterGivenName:         v.RequesterGivenName,
			RequesterSurName:           v.RequesterSurName,
			RequesterUserPrincipalName: v.RequesterUserPrincipalName,
			ApprovalTypeId:             int(v.ApprovalTypeId),
			ApprovalType:               v.ApprovalType,
			ApprovalDescription:        v.ApprovalDescription,
			RequestStatus:              v.RequestStatus,
			ApprovalDate:               v.ApprovalDate,
			OSSsponsor:                 v.OSSsponsor,
			OfferingsAssets:            v.Offeringsassets,
			WillBeCommercialVersion:    v.Willbecommercialversion,
			OSSContributionInformation: v.OSSContributionInformation,
		}

		approvalRequestApprovers, err := GetApprovalRequestApproversByApprovalRequestId(int(v.Id))
		if err != nil {
			return nil, err
		}
		for _, approvalRequestApprover := range approvalRequestApprovers {
			projectApprovalApprover.Approvers = append(projectApprovalApprover.Approvers, approvalRequestApprover.ApproverEmail)
		}

		projectApprovalApprovers = append(projectApprovalApprovers, projectApprovalApprover)
	}

	return projectApprovalApprovers, nil
}
