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
	Coowner                 string
	Description             string
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

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_ProjectApprovals_Select_ById", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func ProjectsSelectByUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Select_ByUserPrincipalName", params)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func ProjectsApprovalUpdateApproverUserPrincipalName(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_ProjectsApproval_Update_ApproverUserPrincipalName", params)
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
		"CoOwner":                 project.Coowner,
		"Description":             project.Description,
		"ConfirmAvaIP":            project.ConfirmAvaIP,
		"ConfirmEnabledSecurity":  project.ConfirmSecIPScan,
		"ConfirmNotClientProject": project.ConfirmNotClientProject,
		"CreatedBy":               user,
		"VisibilityId":            project.Visibility,
		"TFSProjectReference":     project.TFSProjectReference,
	}
	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_Insert", param)
	if err != nil {
		fmt.Println(err)
	}
	id = result[0]["ItemId"].(int64)
	return
}

func DeleteProjectById(id int) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}
	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Delete_ById", param)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func ProjectInsertByImport(param map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Insert", param)
	if err != nil {
		return err
	}
	return nil
}

func ProjectUpdateByImport(param map[string]interface{}) error {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Update_Repo_Info", param)
	if err != nil {
		return err
	}
	return nil
}

func PRProjectsUpdate(project Project, user string) (id int64) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ID":                      project.Id,
		"Name":                    project.Name,
		"CoOwner":                 project.Coowner,
		"Description":             project.Description,
		"ConfirmAvaIP":            project.ConfirmAvaIP,
		"ConfirmEnabledSecurity":  project.ConfirmSecIPScan,
		"ConfirmNotClientProject": project.ConfirmNotClientProject,
		"ModifiedBy":              user,
	}
	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Update", param)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func UpdateProjectEcattIdById(id, ecattId int, user string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":         id,
		"ECATTID":    ecattId,
		"ModifiedBy": user,
	}

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_ECATTID_ById", param)
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
	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_LegalQuestions", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ByRepositorySource", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_IsExisting", param)

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

	result, err := db.ExecuteStoredProcedureWithResult("dbo.PR_Projects_IsExisting_By_GithubId", param)

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

func PopulateProjectsApproval(id int64, email string) (projectApprovals []ProjectApproval) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": id,
		"CreatedBy": email,
	}
	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectsApproval_Populate", param)
	for _, v := range result {
		data := ProjectApproval{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["newcontribution"].(string),
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

func GetFailedProjectApprovalRequests() (projectApprovals []ProjectApproval) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_Failed", nil)

	for _, v := range result {
		data := ProjectApproval{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["newcontribution"].(string),
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

func GetFailedProjectApprovalRequestsObsolete() (projectApprovals []ProjectApproval) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_Failed", nil)

	for _, v := range result {
		data := ProjectApproval{
			Id:                         v["Id"].(int64),
			ProjectId:                  v["ProjectId"].(int64),
			ProjectName:                v["ProjectName"].(string),
			ProjectCoowner:             v["ProjectCoowner"].(string),
			ProjectDescription:         v["ProjectDescription"].(string),
			RequesterGivenName:         v["RequesterGivenName"].(string),
			RequesterSurName:           v["RequesterSurName"].(string),
			RequesterName:              v["RequesterName"].(string),
			RequesterUserPrincipalName: v["RequesterUserPrincipalName"].(string),
			CoownerGivenName:           v["CoownerGivenName"].(string),
			CoownerSurName:             v["CoownerSurName"].(string),
			CoownerName:                v["CoownerName"].(string),
			CoownerUserPrincipalName:   v["CoownerUserPrincipalName"].(string),
			ApprovalTypeId:             v["ApprovalTypeId"].(int64),
			ApprovalType:               v["ApprovalType"].(string),
			ApproverUserPrincipalName:  v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:        v["ApprovalDescription"].(string),
			Newcontribution:            v["newcontribution"].(string),
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

func GetProjectApprovals() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func GetProjectApprovalsByProjectId(id int64) (projectApprovals []ProjectApproval) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_By_ProjectId", param)

	for _, v := range result {
		data := ProjectApproval{
			Id:                        v["Id"].(int64),
			ProjectId:                 v["ProjectId"].(int64),
			ProjectName:               v["ProjectName"].(string),
			ApprovalTypeId:            v["ApprovalTypeId"].(int64),
			ApprovalType:              v["ApprovalType"].(string),
			ApproverUserPrincipalName: v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:       v["ApprovalDescription"].(string),
			RequestStatus:             v["RequestStatus"].(string),
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_By_StatusId", param)
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

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_By_ApprovalSystemGUID", param)

	for _, v := range result {
		projectApproval = ProjectApproval{
			Id:                  v["Id"].(int64),
			ProjectId:           v["ProjectId"].(int64),
			ProjectName:         v["ProjectName"].(string),
			ApprovalTypeId:      v["ApprovalTypeId"].(int64),
			ApprovalType:        v["ApprovalType"].(string),
			ApprovalDescription: v["ApprovalDescription"].(string),
			RequestStatus:       v["RequestStatus"].(string),
		}
	}

	return
}

func UpdateProjectApprovalById(id int, respondedBy string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":          id,
		"RespondedBy": respondedBy,
	}

	_, err := db.ExecuteStoredProcedure("PR_ProjectApprovals_UpdateRespondedBy_ById", param)
	if err != nil {
		return err
	}

	return nil
}

func ProjectsApprovalUpdateGUID(id int64, ApprovalSystemGUID string) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                 id,
		"ApprovalSystemGUID": ApprovalSystemGUID,
	}
	db.ExecuteStoredProcedure("PR_ProjectsApproval_Update_ApprovalSystemGUID", param)
}

func GetProjectForRepoOwner() (repoOwner []RepoOwner) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_ToRepoOwners", nil)

	for _, v := range result {
		data := RepoOwner{
			Id:                v["Id"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		repoOwner = append(repoOwner, data)
	}
	return repoOwner
}

func GetProjectByName(projectName string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Name": projectName,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ByName", param)

	return result
}

func GetProjectById(id int64) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ById", param)

	return result
}

func GetProjectByGithubId(githubId int64) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId": githubId,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ByGithubId", param)

	return result
}

func GetProjectByAssetCode(assetCode string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"AssetCode": assetCode,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_ByAssetCode", param)

	return result
}

func ReposSelectByOffsetAndFilter(offset int, search string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset": offset,
		"Search": search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Repositories_Select_ByOffsetAndFilter", param)
	return result
}

func ReposTotalCountBySearchTerm(search string) int {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"search": search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Repositories_TotalCount_BySearchTerm", param)
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

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_IsArchived_ById", param)
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

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_Visibility_ById", param)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTFSProjectReferenceById(id int64, tFSProjectReference string) error {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                  id,
		"TFSProjectReference": tFSProjectReference,
	}

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_TFSProjectReference_ById", param)
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

	result, err := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_By_DateRange", param)
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

func GetGitHubRepositories() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_Projects_SelectAllGitHub", nil)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SelectReposWithMakePublicRequest() ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	sponsors, err := db.ExecuteStoredProcedureWithResult("PR_Projects_Select_AllWithMakePublicRequest", nil)
	if err != nil {
		return nil, err
	}
	return sponsors, nil
}

func UpdateOssContributionSponsorId(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_OssSponsorIdById", params)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func LegacySearch(params map[string]interface{}) ([]map[string]interface{}, error) {
	db := ConnectDb()
	defer db.Close()

	result, err := db.ExecuteStoredProcedureWithResult("PR_Projects_LegacySearch", params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetProjectApprovalsByUsername(username string) ([]ApprovalRequest, error) {
	db := ConnectDb()
	defer db.Close()

	params := map[string]interface{}{
		"ApproverEmail": username,
	}

	result, err := db.ExecuteStoredProcedureWithResult("PR_ApprovalRequests_Select_ByApproverEmail", params)
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
