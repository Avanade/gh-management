package ghmgmt

import (
	"fmt"
	"strconv"
	"time"
)

type Repository struct {
	Name        string
	Requestor   string
	Description string
}

type ProjectRequest struct {
	Id                         string
	Newcontribution            string
	OSSsponsor                 int
	Offeringsassets            string
	Willbecommercialversion    string
	OSSContributionInformation string
}

type ProjectApproval struct {
	Id                         int64
	ProjectId                  int64
	ProjectName                string
	ProjectCoowner             string
	ProjectDescription         string
	RequesterName              string
	RequesterGivenName         string
	RequesterSurName           string
	RequesterUserPrincipalName string
	CoownerName                string
	CoownerGivenName           string
	CoownerSurName             string
	CoownerUserPrincipalName   string
	ApprovalTypeId             int64
	ApprovalType               string
	ApproverUserPrincipalName  string
	ApprovalDescription        string
	Newcontribution            string
	OSSsponsor                 string
	Offeringsassets            string
	Willbecommercialversion    string
	OSSContributionInformation string
	RequestStatus              string
	ApproveUrl                 string
	RejectUrl                  string
	ApproveText                string
	RejectText                 string
	ApprovalDate               time.Time
}

type Project struct {
	Id                      string
	GithubId                int64
	Name                    string
	AssetCode               string
	Coowner                 string
	Description             string
	Organization            string
	ConfirmAvaIP            bool
	ConfirmSecIPScan        bool
	ConfirmNotClientProject bool
	TFSProjectReference     string
	Visibility              int
}

type ApprovalRequest struct {
	ApprovalSystemGUID         string
	ProjectName                string
	RequestedBy                string
	Description                string
	NewContribution            string
	OSSContributorSponsor      string
	IsAvanadeOfferingAssets    string
	WillBeCommercialVersion    string
	OSSContributionInformation string
	Remarks                    string
	Status                     string
	RespondedBy                string
	ApprovalDate               string
	ApprovalSystemDateSent     string
}

func ProjectApprovalsSelectById(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApproval_Select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func ProjectsSelectByUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_Select_ByUserPrincipalName", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func ProjectsApprovalUpdateApproverUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApproval_Update_ApproverUserPrincipalName", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func PRProjectsInsert(project Project, user string) (id int64) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId":                project.GithubId,
		"Name":                    project.Name,
		"AssetCode":               project.AssetCode,
		"CoOwner":                 project.Coowner,
		"Description":             project.Description,
		"Organization":            project.Organization,
		"ConfirmAvaIP":            project.ConfirmAvaIP,
		"ConfirmEnabledSecurity":  project.ConfirmSecIPScan,
		"ConfirmNotClientProject": project.ConfirmNotClientProject,
		"CreatedBy":               user,
		"VisibilityId":            project.Visibility,
		"TFSProjectReference":     project.TFSProjectReference,
	}
	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_Insert", param)
	if err != nil {
		fmt.Println(err)
	}
	id = result[0]["Id"].(int64)
	return
}

func DeleteProjectById(id int) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositoryId": id,
	}
	_, err := db.ExecuteStoredProcedure("usp_Repository_Delete", param)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ProjectInsertByImport(param map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("usp_Repository_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func ProjectUpdateByImport(param map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("usp_Repository_Update", param)
	if err != nil {
		return err
	}
	return nil
}

func UpdateProjectEcattIdById(id, ecattId int, user string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":         id,
		"ECATTID":    ecattId,
		"ModifiedBy": user,
	}

	_, err := db.ExecuteStoredProcedure("usp_Repository_Update_ECATTID", param)
	if err != nil {
		return err
	}
	return nil
}

func PRProjectsUpdateLegalQuestions(projectRequest ProjectRequest, user string) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                         projectRequest.Id,
		"Newcontribution":            projectRequest.Newcontribution,
		"OSSsponsor":                 projectRequest.OSSsponsor,
		"Avanadeofferingsassets":     projectRequest.Offeringsassets,
		"Willbecommercialversion":    projectRequest.Willbecommercialversion,
		"OSSContributionInformation": projectRequest.OSSContributionInformation,
		"ModifiedBy":                 user,
	}
	_, err := db.ExecuteStoredProcedure("usp_Repository_Update_LegalQuestions", param)
	if err != nil {
		fmt.Println(err)
	}
}

func ProjectsByRepositorySource(repositorySource string) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"RepositorySource": repositorySource,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_Select_ByRepositorySource", param)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ProjectsIsExisting(name string) bool {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"Name": name,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_IsNameExist", param)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if result[0]["Result"] == "1" {
		return true
	} else {
		return false
	}
}

func ProjectsIsExistingByGithubId(githubId int64) bool {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId": githubId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_IsGitHubIdExist", param)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if result[0]["Result"] == "1" {
		return true
	} else {
		return false
	}
}

func GetFailedProjectApprovalRequests() (projectApprovals []ProjectApproval) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("usp_RepositoryApproval_Select_FailedRequest", nil)

	for _, v := range result {
		data := ProjectApproval{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["RepositoryId"].(int64),
			ProjectName:                v["RepositoryName"].(string),
			ProjectDescription:         v["RepositoryName"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApprovalTypeId:             v["RepositoryApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["Newcontribution"].(string),
			OSSsponsor:                 v["OSSsponsor"].(string),
			Offeringsassets:            v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
		}
		projectApprovals = append(projectApprovals, data)
	}

	return
}

func GetProjectApprovalsByProjectId(id int64) (projectApprovals []ProjectApproval) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_RepositoryApproval_Select_ByRecentProjectId", param)

	for _, v := range result {
		data := ProjectApproval{
			Id:                  v["Id"].(int64),
			ProjectId:           v["RepositoryId"].(int64),
			ProjectName:         v["RepositoryName"].(string),
			ApprovalTypeId:      v["RepositoryApprovalTypeId"].(int64),
			ApprovalType:        v["ApprovalType"].(string),
			ApprovalDescription: v["ApprovalDescription"].(string),
			RequestStatus:       v["RequestStatus"].(string),
		}
		if v["ApprovalDate"] != nil {
			data.ApprovalDate = v["ApprovalDate"].(time.Time)
		}

		projectApprovals = append(projectApprovals, data)
	}

	return
}

func GetProjectApprovalsByStatusId(approvalStatusId int64) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalStatusId": approvalStatusId,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_RepositoryApproval_Select_ByApprovalStatusId", param)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetProjectApprovalByGUID(id string) (projectApproval ProjectApproval) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalSystemGUID": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_RepositoryApproval_Select_ByApprovalSystemGUID", param)

	for _, v := range result {
		projectApproval = ProjectApproval{
			Id:                  v["Id"].(int64),
			ProjectId:           v["RepositoryId"].(int64),
			ProjectName:         v["RepositoryName"].(string),
			ApprovalTypeId:      v["RepositoryApprovalTypeId"].(int64),
			ApprovalType:        v["ApprovalType"].(string),
			ApprovalDescription: v["ApprovalDescription"].(string),
			RequestStatus:       v["RequestStatus"].(string),
		}
	}

	return
}

func ProjectsApprovalUpdateGUID(id int64, ApprovalSystemGUID string) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                 id,
		"ApprovalSystemGUID": ApprovalSystemGUID,
	}
	db.ExecuteStoredProcedure("usp_RepositoryApproval_Update_ApprovalSystemGUID", param)
}

func GetProjectIdByOrgName(orgName, repoName string) (int64, error) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"AssetCode":    repoName,
		"Organization": orgName,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_GetProjectId_ByAssetCodeAndOrganization", param)
	if err != nil {
		return 0, err
	}

	if result == nil {
		return 0, fmt.Errorf("project with the organization name '%v' and the repository name '%v' does not exist", orgName, repoName)
	}

	return result[0]["Id"].(int64), nil
}

func GetProjectById(id int64) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_Repository_Select_ById", param)

	return result
}

func GetProjectByGithubId(githubId int64) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId": githubId,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_Repository_Select_ByGitHubId", param)

	return result
}

func GetProjectByAssetCode(assetCode string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"AssetCode": assetCode,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_Repository_Select_ByAssetCode", param)

	return result
}

func ReposSelectByOffsetAndFilter(offset int, search string, filterType int, filter string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset":     offset,
		"Search":     search,
		"FilterType": filterType,
		"Filter":     filter,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_Repository_Select_ByOption", param)
	return result
}

func ReposTotalCountBySearchTerm(search string, filterType int, filter string) int {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"FilterType": filterType,
		"Filter":     filter,
		"Search":     search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("usp_Repository_TotalCount_ByOption", param)
	if result == nil {
		return 0
	}
	total, err := strconv.Atoi(fmt.Sprint(result[0]["Total"]))
	if err != nil {
		return 0
	}
	return total
}

func UpdateProjectIsArchived(id int64, isArchived bool) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":         id,
		"IsArchived": isArchived,
	}

	_, err := db.ExecuteStoredProcedure("usp_Repository_Update_IsArchived", param)
	if err != nil {
		return err
	}

	return nil
}

func UpdateProjectVisibilityId(id int64, visibilityId int64) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":           id,
		"VisibilityId": visibilityId,
	}

	_, err := db.ExecuteStoredProcedure("usp_Repository_Update_Visibility", param)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTFSProjectReferenceById(id int64, tFSProjectReference, organization string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                  id,
		"TFSProjectReference": tFSProjectReference,
		"Organization":        organization,
	}

	_, err := db.ExecuteStoredProcedure("usp_Repository_Update_TFSProjectReference", param)
	if err != nil {
		return err
	}

	return nil
}

func GetRequestedReposByDateRange(start time.Time, end time.Time) ([]Repository, error) {
	var requestedRepos []Repository
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Start": start.Format("2006-01-02"),
		"End":   end.Format("2006-01-02"),
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_Select_ByDateRange", param)
	if err != nil {
		return nil, err
	}

	for _, v := range result {
		data := Repository{
			Name:        v["Name"].(string),
			Requestor:   v["CreatedBy"].(string),
			Description: v["Description"].(string),
		}
		requestedRepos = append(requestedRepos, data)
	}

	return requestedRepos, nil
}

func LegacySearch(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("usp_Repository_LegacySearch", params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetProjectApprovalsByUsername(username string) ([]ApprovalRequest, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"ApproverUserPrincipalName": username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("usp_ApprovalRequest_Select_ByApproverUserPrincipalName", params)
	if err != nil {
		return nil, err
	}

	var approvalRequests []ApprovalRequest

	for _, v := range result {
		approvalRequest := ApprovalRequest{
			ApprovalSystemGUID:         fmt.Sprintf("%v", v["ApprovalSystemGUID"]),
			ProjectName:                fmt.Sprintf("%v", v["ProjectName"]),
			RequestedBy:                fmt.Sprintf("%v", v["RequestedBy"]),
			Description:                fmt.Sprintf("%v", v["Description"]),
			NewContribution:            fmt.Sprintf("%v", v["NewContribution"]),
			OSSContributorSponsor:      fmt.Sprintf("%v", v["OSSContributionSponsor"]),
			IsAvanadeOfferingAssets:    fmt.Sprintf("%v", v["IsAvanadeOfferingAssets"]),
			WillBeCommercialVersion:    fmt.Sprintf("%v", v["WillBeCommercialVersion"]),
			OSSContributionInformation: fmt.Sprintf("%v", v["OSSContributionInformation"]),
			Remarks:                    fmt.Sprintf("%v", v["Remarks"]),
			Status:                     fmt.Sprintf("%v", v["Status"]),
			RespondedBy:                fmt.Sprintf("%v", v["RespondedBy"]),
			ApprovalDate:               fmt.Sprintf("%v", v["ApprovalDate"]),
			ApprovalSystemDateSent:     fmt.Sprintf("%v", v["ApprovalSystemDateSent"]),
		}

		approvalRequests = append(approvalRequests, approvalRequest)
	}

	return approvalRequests, err
}
