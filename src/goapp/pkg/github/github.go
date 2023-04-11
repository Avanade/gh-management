package githubAPI

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"main/models"
	"main/pkg/email"
	"main/pkg/envvar"
	ghmgmt "main/pkg/ghmgmtdb"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func createClient(token string) *github.Client {
	// create github oauth client from token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func CreatePrivateGitHubRepository(data models.TypNewProjectReqBody, requestor string) (*github.Repository, error) {
	client := createClient(os.Getenv("GH_TOKEN"))
	owner := os.Getenv("GH_ORG_INNERSOURCE")
	repoRequest := &github.TemplateRepoRequest{
		Name:        &data.Name,
		Owner:       &owner,
		Description: &data.Description,
		Private:     github.Bool(true),
	}

	repo, _, err := client.Repositories.CreateFromTemplate(context.Background(), os.Getenv("GH_REPO_TEMPLATE"), os.Getenv("GH_REPO_TEMPLATE_NAME"), repoRequest)
	if err != nil {
		return nil, err
	}

	_, err = AddCollaborator(data, requestor)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func CreateInternalGitHubRepository(data models.TypNewProjectReqBody, requestor string) (*github.Repository, error) {
	// Set your access token or personal access token
	authToken := os.Getenv("GH_TOKEN")

	// Define the repository object with the template repository name
	repoData := map[string]interface{}{
		"name":                &data.Name,
		"private":             false,
		"visibility":          "private",
		"template_repository": os.Getenv("GH_REPO_TEMPLATE_NAME"),
	}

	// Convert the repository object to a JSON string
	jsonData, err := json.Marshal(repoData)
	if err != nil {
		return nil, err
	}

	// Send the POST request to create the repository
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/orgs/%s/repos", os.Getenv("GH_ORG_INNERSOURCE")), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code to make sure the repository was created successfully
	if resp.StatusCode == 201 {
		// Decode the response JSON into a Repository struct
		var repo *github.Repository
		err = json.NewDecoder(resp.Body).Decode(&repo)
		if err != nil {
			return nil, err
		}
		return repo, nil
	} else {
		return nil, errors.New(resp.Status)
	}
}

type Enterprise struct {
	IsEnterprise bool `json:"is_enterprise"`
}

func GetIsOrgEnterprise() (bool, error) {
	// Set your access token or personal access token
	authToken := os.Getenv("GH_TOKEN")

	// Set the name of the organization you want to check
	orgName := os.Getenv("GH_ORG_INNERSOURCE")

	// Send the GET request to the enterprise endpoint
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/enterprises/%s", orgName), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Decode the response JSON into an Enterprise struct
	var enterprise Enterprise
	err = json.NewDecoder(resp.Body).Decode(&enterprise)
	if err != nil {
		return false, err
	}

	return enterprise.IsEnterprise, nil
}

func AddCollaborator(data models.TypNewProjectReqBody, requestor string) (*github.Response, error) {
	client := createClient(os.Getenv("GH_TOKEN"))
	owner := os.Getenv("GH_ORG_INNERSOURCE")
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: "admin",
	}
	if data.Coowner != requestor {
		GHUser := ghmgmt.Users_Get_GHUser(requestor)
		_, _, err := client.Repositories.AddCollaborator(context.Background(), owner, data.Name, GHUser, opts)
		if err != nil {
			return nil, err
		}
	}
	GHUser := ghmgmt.Users_Get_GHUser(data.Coowner)
	_, resp, err := client.Repositories.AddCollaborator(context.Background(), owner, data.Name, GHUser, opts)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetRepository(repoName string, org string) (*github.Repository, error) {
	client := createClient(os.Getenv("GH_TOKEN"))
	owner := org
	repo, _, err := client.Repositories.Get(context.Background(), owner, repoName)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func IsArchived(repoName string, org string) (bool, error) {
	repo, err := GetRepository(repoName, org)
	if err != nil {
		return false, err
	}

	return repo.GetArchived(), nil
}

func Repo_IsExisting(repoName string) (bool, error) {
	exists := false
	organizations := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	for _, org := range organizations {
		_, err := GetRepository(repoName, org)
		if err != nil {
			if strings.Contains(err.Error(), "Not Found") {
				continue
			} else {
				return false, err
			}
		} else {
			exists = true
		}
	}

	return exists, nil
}

func GetRepositoriesFromOrganization(org string) ([]Repo, error) {
	client := createClient(os.Getenv("GH_TOKEN"))
	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{Type: "all", Sort: "full_name", ListOptions: github.ListOptions{PerPage: 30}}

	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opt)
		if err != nil {
			if resp.Response.StatusCode == 403 {
				return nil, nil
			} else {
				return nil, err
			}
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	var repoList []Repo
	for _, repo := range allRepos {
		r := Repo{
			GithubId:    repo.GetID(),
			FullName:    repo.GetFullName(),
			Name:        repo.GetName(),
			Link:        repo.GetHTMLURL(),
			Org:         org,
			Description: repo.GetDescription(),
			Private:     repo.GetPrivate(),
			Created:     repo.GetCreatedAt(),
			IsArchived:  repo.GetArchived(),
			Visibility:  repo.GetVisibility(),
		}
		repoList = append(repoList, r)
	}

	return repoList, nil
}

func SetProjectVisibility(projectName string, visibility string, org string) error {
	client := &http.Client{}
	urlPath := fmt.Sprintf("https://api.github.com/repos/%s/%s", org, projectName)
	postBody, _ := json.Marshal(map[string]string{
		"visibility": visibility,
	})
	reqBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPatch, urlPath, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+envvar.GetEnvVar("GH_TOKEN", ""))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusUnprocessableEntity {
		return errors.New("Failed to make repository " + visibility)
	}

	return nil
}

func ArchiveProject(projectName string, archive bool, org string) error {
	client := &http.Client{}
	urlPath := fmt.Sprintf("https://api.github.com/repos/%s/%s", org, projectName)
	postBody, _ := json.Marshal(map[string]bool{
		"archived": archive,
	})
	reqBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest(http.MethodPatch, urlPath, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+envvar.GetEnvVar("GH_TOKEN", ""))

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func TransferRepository(repo string, owner string, newOwner string) error {
	client := createClient(os.Getenv("GH_TOKEN"))
	opt := github.TransferRequest{NewOwner: newOwner}

	_, _, err := client.Repositories.Transfer(context.Background(), owner, repo, opt)
	if err != nil {
		return err
	}
	return nil
}

type Repo struct {
	GithubId    int64            `json:"id"`
	FullName    string           `json:"repoFullName"`
	Name        string           `json:"repoName"`
	Link        string           `json:"repoLink"`
	Org         string           `json:"org"`
	Description string           `json:"description"`
	Private     bool             `json:"private"`
	Created     github.Timestamp `json:"created"`
	IsArchived  bool             `json:"archived"`
	Visibility  string           `json:"visibility"`
}

func OrganizationsIsMember(token string, GHUser string) (bool, bool, error) {
	client := createClient(token)
	OrgInnerSource := os.Getenv("GH_ORG_INNERSOURCE")
	OrgInnerSourceIsMember, _, err := client.Organizations.IsMember(context.Background(), OrgInnerSource, GHUser)

	OrgOuterSource := os.Getenv("GH_ORG_OPENSOURCE")
	OrgOuterSourceIsMember, _, err := client.Organizations.IsMember(context.Background(), OrgOuterSource, GHUser)
	return OrgInnerSourceIsMember, OrgOuterSourceIsMember, err
}

func OrganizationInvitation(token string, username string, org string) *github.Invitation {
	client := createClient(token)
	Email := ""
	Role := "direct_member"
	teamid := []int64{}
	user, _, _ := client.Users.Get(context.Background(), username)
	intid2 := user.ID
	options := *&github.CreateOrgInvitationOptions{InviteeID: intid2, Email: &Email, Role: &Role, TeamID: teamid}

	invite, _, _ := client.Organizations.CreateOrgInvitation(context.Background(), org, &options)

	return invite
}

func ListOutsideCollaborators(token string, org string) []*github.User {
	client := createClient(token)

	options := *&github.ListOutsideCollaboratorsOptions{}

	collabs, _, _ := client.Organizations.ListOutsideCollaborators(context.Background(), org, &options)

	return collabs
}
func RemoveOutsideCollaborator(token string, org string, username string) *github.Response {
	client := createClient(token)

	repons, err := client.Organizations.RemoveOutsideCollaborator(context.Background(), org, username)

	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	return repons
}
func RemoveOrganizationsMember(token string, org string, username string) *github.Response {
	client := createClient(token)

	repons, err := client.Organizations.RemoveMember(context.Background(), org, username)

	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	return repons
}
func RepositoriesListCollaborators(token string, org string, repo string) []*github.User {
	client := createClient(token)
	options := *&github.ListCollaboratorsOptions{}
	ListCollabs, _, err := client.Repositories.ListCollaborators(context.Background(), org, repo, &options)

	if err != nil {

		fmt.Println(err)
	}

	return ListCollabs
}
func OrgListMembers(token string, org string) []*github.User {
	client := createClient(token)

	ListCollabs, _, err := client.Organizations.ListMembers(context.Background(), org, nil)

	if err != nil {

		fmt.Println(err)
	}

	return ListCollabs
}

func EmailAdmin(admin string, adminemail string, reponame string, outisideCollab []string) {
	e := time.Now()

	link := "https://github.com/" + os.Getenv("GH_ORG_OPENSOURCE") + "/" + reponame
	link = "<a href=\"" + link + "\">" + reponame + "</a>"
	Collablist := "</p> <table  >"
	for _, collab := range outisideCollab {
		Collablist = Collablist + " <tr> <td>" + collab + " </td></tr>"
	}
	Collablist = Collablist + " </table  > <p>"
	body := fmt.Sprintf("<p>Hello %s ,  </p>  \n<p>This is to inform you that your Github repository <b> %s </b> has %o outside collaborator/s. </p> %s  This email was sent to the admins of the repository.  </p> \n <p>OSPO</p>", admin, link, len(outisideCollab), Collablist)

	m := email.TypEmailMessage{
		Subject: "GitHub Repo Collaborators Scan",
		Body:    body,
		To:      adminemail,
	}

	email.SendEmail(m)
	fmt.Printf(" GitHub Repo Collaborators Scan on %s was sent.", e)
}
