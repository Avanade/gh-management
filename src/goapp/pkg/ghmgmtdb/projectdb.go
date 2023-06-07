package ghmgmt

import (
	"fmt"
	"main/models"
	"strconv"
	"time"
)

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

func PRProjectsInsert(body models.TypNewProjectReqBody, user string) (id int64) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId":                body.GithubId,
		"Name":                    body.Name,
		"CoOwner":                 body.Coowner,
		"Description":             body.Description,
		"ConfirmAvaIP":            body.ConfirmAvaIP,
		"ConfirmEnabledSecurity":  body.ConfirmSecIPScan,
		"ConfirmNotClientProject": body.ConfirmNotClientProject,
		"CreatedBy":               user,
		"VisibilityId":            body.Visibility,
		"TFSProjectReference":     body.TFSProjectReference,
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

func PRProjectsUpdate(body models.TypNewProjectReqBody, user string) (id int64) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ID":                      body.Id,
		"Name":                    body.Name,
		"CoOwner":                 body.Coowner,
		"Description":             body.Description,
		"ConfirmAvaIP":            body.ConfirmAvaIP,
		"ConfirmEnabledSecurity":  body.ConfirmSecIPScan,
		"ConfirmNotClientProject": body.ConfirmNotClientProject,
		"ModifiedBy":              user,
	}
	_, err := db.ExecuteStoredProcedure("dbo.PR_Projects_Update", param)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func PRProjectsUpdateLegalQuestions(body models.TypeMakeProjectPublicReqBody, user string) {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id":                         body.Id,
		"Newcontribution":            body.Newcontribution,
		"OSSsponsor":                 body.OSSsponsor,
		"Avanadeofferingsassets":     body.Avanadeofferingsassets,
		"Willbecommercialversion":    body.Willbecommercialversion,
		"OSSContributionInformation": body.OSSContributionInformation,
		"ModifiedBy":                 user,
	}
	_, err := db.ExecuteStoredProcedure("PR_Projects_Update_LegalQuestions", param)
	if err != nil {
		fmt.Println(err)
	}
}

func Projects_ByRepositorySource(repositorySource string) ([]map[string]interface{}, error) {
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

func Projects_IsExisting(body models.TypNewProjectReqBody) bool {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{

		"Name": body.Name,
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

func Projects_IsExisting_By_GithubId(body models.TypNewProjectReqBody) bool {

	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"GithubId": body.GithubId,
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

func PopulateProjectsApproval(id int64, email string) (ProjectApprovals []models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ProjectId": id,
		"CreatedBy": email,
	}
	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectsApproval_Populate", param)
	for _, v := range result {
		data := models.TypProjectApprovals{
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
			Avanadeofferingsassets:     v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
		}
		ProjectApprovals = append(ProjectApprovals, data)
	}

	return
}

func GetFailedProjectApprovalRequests() (ProjectApprovals []models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_Failed", nil)

	for _, v := range result {
		data := models.TypProjectApprovals{
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
			Avanadeofferingsassets:     v["Avanadeofferingsassets"].(string),
			Willbecommercialversion:    v["Willbecommercialversion"].(string),
			OSSContributionInformation: v["OSSContributionInformation"].(string),
			RequestStatus:              v["RequestStatus"].(string),
		}
		ProjectApprovals = append(ProjectApprovals, data)
	}

	return
}

func GetProjectApprovalsByProjectId(id int64) (ProjectApprovals []models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Id": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_By_ProjectId", param)

	for _, v := range result {
		data := models.TypProjectApprovals{
			Id:                        v["Id"].(int64),
			ProjectId:                 v["ProjectId"].(int64),
			ProjectName:               v["ProjectName"].(string),
			ApprovalTypeId:            v["ApprovalTypeId"].(int64),
			ApprovalType:              v["ApprovalType"].(string),
			ApproverUserPrincipalName: v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:       v["ApprovalDescription"].(string),
			RequestStatus:             v["RequestStatus"].(string),
		}
		ProjectApprovals = append(ProjectApprovals, data)
	}

	return
}

func GetProjectApprovalByGUID(id string) (ProjectApproval models.TypProjectApprovals) {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"ApprovalSystemGUID": id,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_ProjectApprovals_Select_By_ApprovalSystemGUID", param)

	for _, v := range result {
		data := models.TypProjectApprovals{
			Id:                        v["Id"].(int64),
			ProjectId:                 v["ProjectId"].(int64),
			ProjectName:               v["ProjectName"].(string),
			ApprovalTypeId:            v["ApprovalTypeId"].(int64),
			ApprovalType:              v["ApprovalType"].(string),
			ApproverUserPrincipalName: v["ApproverUserPrincipalName"].(string),
			ApprovalDescription:       v["ApprovalDescription"].(string),
			RequestStatus:             v["RequestStatus"].(string),
		}
		ProjectApproval = data
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
	db.ExecuteStoredProcedure("PR_ProjectsApproval_Update_ApprovalSystemGUID", param)
}

func GetProjectForRepoOwner() (RepoOwner []models.TypRepoOwner) {
	db := ConnectDb()
	defer db.Close()

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Projects_ToRepoOwners", nil)

	for _, v := range result {
		data := models.TypRepoOwner{
			Id:                v["Id"].(int64),
			UserPrincipalName: v["UserPrincipalName"].(string),
		}
		RepoOwner = append(RepoOwner, data)
	}
	return RepoOwner
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

func Repos_Select_ByOffsetAndFilter(offset int, search string) []map[string]interface{} {
	db := ConnectDb()
	defer db.Close()

	param := map[string]interface{}{
		"Offset": offset,
		"Search": search,
	}

	result, _ := db.ExecuteStoredProcedureWithResult("PR_Repositories_Select_ByOffsetAndFilter", param)
	return result
}

func Repos_TotalCount_BySearchTerm(search string) int {
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

func GetRequestedReposByDateRange(start time.Time, end time.Time) ([]models.TypBasicRepo, error) {
	var RequestedRepos []models.TypBasicRepo
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
		data := models.TypBasicRepo{
			Name:        v["Name"].(string),
			Requestor:   v["CreatedBy"].(string),
			Description: v["Description"].(string),
		}
		RequestedRepos = append(RequestedRepos, data)
	}

	return RequestedRepos, nil
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
