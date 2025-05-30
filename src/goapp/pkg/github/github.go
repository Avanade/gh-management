package githubAPI

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/shurcooL/githubv4"

	db "main/pkg/ghmgmtdb"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type Repo struct {
	GithubId            int64            `json:"id"`
	FullName            string           `json:"repoFullName"`
	Name                string           `json:"repoName"`
	Link                string           `json:"repoLink"`
	Org                 string           `json:"org"`
	Description         string           `json:"description"`
	Private             bool             `json:"private"`
	Created             github.Timestamp `json:"created"`
	IsArchived          bool             `json:"archived"`
	Visibility          string           `json:"visibility"`
	TFSProjectReference string
	Topics              []string
	Role                string `json:"roleName"`
}

func CreateClient(token string) *github.Client {
	// create github oauth client from token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func CreatePrivateGitHubRepository(name, description, requestor string) (*github.Repository, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	owner := os.Getenv("GH_ORG_INNERSOURCE")
	repoRequest := &github.TemplateRepoRequest{
		Name:        &name,
		Owner:       &owner,
		Description: &description,
		Private:     github.Bool(true),
	}

	repo, _, err := client.Repositories.CreateFromTemplate(context.Background(), os.Getenv("GH_REPO_TEMPLATE"), os.Getenv("GH_REPO_TEMPLATE_NAME"), repoRequest)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func IsEnterpriseOrg() (bool, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	orgName := os.Getenv("GH_ORG_INNERSOURCE")
	org, _, err := client.Organizations.Get(context.Background(), orgName)
	if err != nil {
		return false, err
	}
	return *org.Plan.Name == "enterprise", err
}

func AddCollaborator(owner string, repo string, user string, permission string) (*github.Response, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: permission,
	}

	_, resp, err := client.Repositories.AddCollaborator(context.Background(), owner, repo, user, opts)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func RemoveCollaborator(owner string, repo string, user string, permission string) (*github.Response, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))

	resp, err := client.Repositories.RemoveCollaborator(context.Background(), owner, repo, user)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetRepository(repoName string, org string) (*github.Repository, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	owner := org
	repo, _, err := client.Repositories.Get(context.Background(), owner, repoName)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func GetPermissionLevel(repoOwner string, repoName string, username string) (string, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	permission, _, err := client.Repositories.GetPermissionLevel(context.Background(), repoOwner, repoName, username)
	if err != nil {
		return "", err
	}
	return permission.User.GetRoleName(), nil
}

func GetRepositoryReadmeById(owner, repoName string) (string, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))

	readme, resp, err := client.Repositories.GetReadme(context.Background(), owner, repoName, nil)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("unexpected response status: %v", resp.Status)
	}

	// Decode the base64-encoded content of the README
	decodedContent, err := base64.StdEncoding.DecodeString(*readme.Content)
	if err != nil {
		log.Fatal(err)
	}

	return string(decodedContent), nil
}

func IsArchived(repoName string, org string) (bool, error) {
	repo, err := GetRepository(repoName, org)
	if err != nil {
		return false, err
	}

	return repo.GetArchived(), nil
}

func IsRepoExisting(repoName string) (bool, error) {
	exists := false
	orgs := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	isEnabled := db.NullBool{Value: true}
	regOrgs, _ := db.SelectRegionalOrganization(&isEnabled)

	for _, regOrg := range regOrgs {
		if !regOrg.IsIndexRepoEnabled {
			continue
		}
		orgs = append(orgs, regOrg.Name)
	}

	for _, org := range orgs {
		repo, err := GetRepository(repoName, org)
		if err != nil {
			if strings.Contains(err.Error(), "Not Found") {
				continue
			} else {
				return false, err
			}
		} else {
			if *repo.GetOwner().Login == org {
				exists = true
			}
		}
	}

	return exists, nil
}

func GetCollaboratorRepositoriesFromOrganization(token, org, user string) ([]Repo, error) {
	repos, err := GetRepositoriesFromOrganization(org)
	if err != nil {
		return nil, err
	}

	var collaboratorRepositories []Repo

	for _, repo := range repos {

		// Successfully fetches the collaborators of the repository
		collaborators, err := GetRepositoryDirectCollaborators(token, org, repo.Name)
		if err != nil {
			log.Printf("Error checking if user %s is a collaborator for repository %s: %v", user, repo.Name, err)
			continue
		}

		for _, collaborator := range collaborators {

			if user == *collaborator.Login {
				collaboratorRepositories = append(collaboratorRepositories, repo)
				break
			}
		}

		if len(collaborators) == 0 {
			log.Printf("User %s is not a collaborator for repository %s in organization %s", user, repo.Name, org)
		}
	}
	return collaboratorRepositories, nil
}

func GetRepositoriesFromOrganization(org string) ([]Repo, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
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
			GithubId:            repo.GetID(),
			FullName:            repo.GetFullName(),
			Name:                repo.GetName(),
			Link:                repo.GetHTMLURL(),
			Org:                 *repo.GetOwner().Login,
			Description:         repo.GetDescription(),
			Private:             repo.GetPrivate(),
			Created:             repo.GetCreatedAt(),
			IsArchived:          repo.GetArchived(),
			Visibility:          repo.GetVisibility(),
			TFSProjectReference: repo.GetHTMLURL(),
			Topics:              repo.Topics,
			Role:                repo.GetRoleName(),
		}
		repoList = append(repoList, r)
	}

	return repoList, nil
}

func SetProjectVisibility(projectName string, visibility string, org string) (*github.Response, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	opt := &github.Repository{Visibility: github.String(visibility)}

	_, resp, err := client.Repositories.Edit(context.Background(), org, projectName, opt)
	if err != nil {
		return resp, err

	}
	return resp, nil
}

func ArchiveProject(projectName string, archive bool, org string) error {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	opt := &github.Repository{Archived: github.Bool(archive)}

	_, _, err := client.Repositories.Edit(context.Background(), org, projectName, opt)
	if err != nil {
		return err
	}
	return nil
}

func TransferRepository(name string, owner string, newOwner string) (*github.Repository, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	opt := github.TransferRequest{NewOwner: newOwner}

	repo, resp, err := client.Repositories.Transfer(context.Background(), owner, name, opt)
	if resp.StatusCode != 202 {
		return nil, err
	}
	return repo, nil
}

func IsOrganizationMember(token, org, ghUser string) (bool, error) {
	client := CreateClient(token)
	isOrgMember, _, err := client.Organizations.IsMember(context.Background(), org, ghUser)
	return isOrgMember, err
}

func IsRepositoryCollaborator(token, owner, repo, ghUser string) (bool, error) {
	client := CreateClient(token)
	isRepoMember, _, err := client.Repositories.IsCollaborator(context.Background(), owner, repo, ghUser)
	return isRepoMember, err
}

func GetRepositoryDirectCollaborators(token, owner, repo string) ([]*github.User, error) {
	client := CreateClient(token)
	opts := &github.ListCollaboratorsOptions{Affiliation: "direct"}
	collaborators, ghResponse, _ := client.Repositories.ListCollaborators(context.Background(), owner, repo, opts)

	if ghResponse.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected response status: %v", ghResponse.Status)
	}
	return collaborators, nil
}

func UserMembership(token, org, ghUser string) (*github.Membership, error) {
	client := CreateClient(token)
	membership, _, err := client.Organizations.GetOrgMembership(context.Background(), ghUser, org)
	return membership, err
}

func OrganizationInvitation(token string, username string, org string) *github.Invitation {
	client := CreateClient(token)
	REINSTATE_ROLE := "reinstate"

	teamid := []int64{}
	email := ""
	user, _, _ := client.Users.Get(context.Background(), username)
	options := &github.CreateOrgInvitationOptions{InviteeID: user.ID, Email: &email, Role: &REINSTATE_ROLE, TeamID: teamid}

	var invite *github.Invitation

	invite, resp, _ := client.Organizations.CreateOrgInvitation(context.Background(), org, options)

	if resp.StatusCode != 201 {
		DIRECT_MEMBER_ROLE := "direct_member"
		options.Role = &DIRECT_MEMBER_ROLE
		invite, _, _ = client.Organizations.CreateOrgInvitation(context.Background(), org, options)
	}

	return invite
}

func ListPendingOrgInvitations(token, org string) []*github.Invitation {
	client := CreateClient(token)
	options := &github.ListOptions{PerPage: 30}

	var allPendingInvitations []*github.Invitation

	for {
		pendingInvitations, resp, err := client.Organizations.ListPendingOrgInvitations(context.Background(), org, options)
		if err != nil {
			log.Printf("ERROR : %s", err.Error())
			return nil
		}

		allPendingInvitations = append(allPendingInvitations, pendingInvitations...)
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return allPendingInvitations
}

func ListOutsideCollaborators(token string, org string) []*github.User {
	client := CreateClient(token)

	options := &github.ListOutsideCollaboratorsOptions{ListOptions: github.ListOptions{PerPage: 30}}

	var collaborators []*github.User
	for {
		collabs, resp, err := client.Organizations.ListOutsideCollaborators(context.Background(), org, options)
		if err != nil {
			log.Printf("ERROR : %s", err.Error())
			return nil
		}

		collaborators = append(collaborators, collabs...)
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	return collaborators
}

func RemoveOutsideCollaborator(token string, org string, username string) *github.Response {
	client := CreateClient(token)

	repons, err := client.Organizations.RemoveOutsideCollaborator(context.Background(), org, username)

	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	return repons
}

func ConvertMemberToOutsideCollaborator(token string, org string, username string) *github.Response {
	client := CreateClient(token)

	repons, err := client.Organizations.ConvertMemberToOutsideCollaborator(context.Background(), org, username)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	return repons
}

func RemoveOrganizationsMember(token string, org string, username string) *github.Response {
	client := CreateClient(token)

	repons, err := client.Organizations.RemoveMember(context.Background(), org, username)

	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	return repons
}

func RepositoriesListCollaborators(token string, org string, repo string, role string, affiliations string) []*github.User {
	client := CreateClient(token)
	options := &github.ListCollaboratorsOptions{Permission: role, Affiliation: affiliations, ListOptions: github.ListOptions{PerPage: 30}}

	var collaborators []*github.User
	for {
		listCollaborators, resp, err := client.Repositories.ListCollaborators(context.Background(), org, repo, options)
		if err != nil {
			log.Printf("ERROR : %s", err.Error())
			return nil
		}

		collaborators = append(collaborators, listCollaborators...)
		if resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}
	return collaborators
}

func OrgListMembers(token string, org string, role string) ([]*github.User, error) {
	client := CreateClient(token)

	opts := &github.ListMembersOptions{Role: role, ListOptions: github.ListOptions{PerPage: 30}}

	var members []*github.User

	for {
		listMembers, resp, err := client.Organizations.ListMembers(context.Background(), org, opts)
		if err != nil {
			log.Printf("ERROR : %s", err.Error())
			return nil, err
		}

		members = append(members, listMembers...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return members, nil
}

func GetOrganizations(token string) ([]*github.Organization, error) {
	client := CreateClient(token)
	opts := &github.ListOptions{PerPage: 30}
	var orgs []*github.Organization

	for {
		listOrgs, resp, err := client.Organizations.List(context.Background(), "", opts)
		if err != nil {
			log.Printf("ERROR : %s", err.Error())
			return nil, err
		}

		orgs = append(orgs, listOrgs...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return orgs, nil

}

func GetTeam(token string, org string, slug string) (*github.Team, error) {
	client := CreateClient(token)

	team, resp, err := client.Teams.GetTeamBySlug(context.Background(), org, slug)
	if err != nil && resp.StatusCode != 404 {
		log.Printf("ERROR : %s", err.Error())
		return nil, err
	}

	return team, nil
}

func CreateTeam(token string, org string, teamName string) (*github.Team, error) {
	client := CreateClient(token)
	newTeam := github.NewTeam{
		Name: teamName,
	}

	team, _, err := client.Teams.CreateTeam(context.Background(), org, newTeam)
	if err != nil {
		log.Printf("ERROR : %s", err.Error())
		return nil, err
	}

	return team, nil
}

func AddMemberToTeam(token string, org string, slug string, user string, role string) (*github.Membership, error) {
	client := CreateClient(token)
	opts := &github.TeamAddTeamMembershipOptions{Role: role}

	teamMembership, _, err := client.Teams.AddTeamMembershipBySlug(context.Background(), org, slug, user, opts)
	if err != nil {
		log.Printf("ERROR : %s", err.Error())
		return nil, err
	}

	return teamMembership, nil
}

func RemoveEnterpriseMember(token string, enterpriseId string, userId string) error {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	var mutation struct {
		RemoveEnterpriseMember struct {
			ClientMutationId string
		} `graphql:"removeEnterpriseMember(input: $input)"`
	}

	input := githubv4.RemoveEnterpriseMemberInput{
		EnterpriseID: githubv4.ID(enterpriseId),
		UserID:       githubv4.ID(userId),
	}

	err := client.Mutate(context.Background(), &mutation, input, nil)
	if err != nil {
		return err
	}

	log.Printf("User %s removed from enterprise %s", userId, enterpriseId)
	return nil
}

func GetOrganizationsByGitHubName(username string, token string) (*GetOrganizationsByGithubNameResult, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var result GetOrganizationsByGithubNameResult
	duplicatedOrgs := make(map[string]bool)

	for {
		var queryResult GetOrganizationsByGitHubNameQuery
		variables := map[string]interface{}{
			"login": githubv4.String(username),
		}
		err := client.Query(context.Background(), &queryResult, variables)
		if err != nil {
			return nil, err
		}

		if queryResult.User.Organizations.Edges == nil {
			return nil, fmt.Errorf("no organizations found for user: %s", username)
		}

		for _, org := range queryResult.User.Organizations.Edges {
			if duplicatedOrgs[string(org.Node.Login)] {
				return &result, nil
			}

			duplicatedOrgs[string(org.Node.Login)] = true
			result.Organizations = append(result.Organizations, Organization{
				Login:      string(org.Node.Login),
				DatabaseId: int64(org.Node.DatabaseId),
			})
		}

		if len(queryResult.User.Organizations.Edges) == 0 {
			break
		}

		isEnabled := db.NullBool{Value: true}
		regOrgs, _ := db.SelectRegionalOrganization(&isEnabled)
		orgs := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}
		for _, org := range regOrgs {
			if org.IsRegionalOrganization {
				orgs = append(orgs, org.Name)
			}
		}
	}
	return &result, nil
}

func GetOrganizationsWithinEnterprise(enterprise string, token string) (*GetOrganizationsWithinEnterpriseResult, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var result GetOrganizationsWithinEnterpriseResult
	var cursor *githubv4.String

	for {
		var queryResult GetOrganizationsWithinEnterpriseQuery
		variables := map[string]interface{}{
			"enterprise": githubv4.String(enterprise),
			"cursor":     cursor,
		}
		err := client.Query(context.Background(), &queryResult, variables)
		if err != nil {
			return nil, err
		}

		for _, org := range queryResult.Enterprise.Organization.Nodes {
			result.Organizations = append(result.Organizations, Organization{
				Login:      string(org.Login),
				DatabaseId: int64(org.DatabaseId),
			})
		}

		if !queryResult.Enterprise.Organization.PageInfo.HasNextPage {
			break
		}

		cursor = &queryResult.Enterprise.Organization.PageInfo.EndCursor
	}

	return &result, nil
}

type customTransport struct {
	Transport http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-GitHub-Next-Global-ID", "1")
	return t.Transport.RoundTrip(req)
}

func GetMembersByEnterprise(enterprise string, token string) (*GetMembersByEnterpriseResult, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	httpClient.Transport = &customTransport{Transport: httpClient.Transport}
	client := githubv4.NewClient(httpClient)

	var result GetMembersByEnterpriseResult
	var after *githubv4.String

	start := time.Now()
	for {
		var queryResult GetMembersByEnterpriseQuery
		variables := map[string]interface{}{
			"enterprise": githubv4.String(enterprise),
			"after":      after,
		}
		err := client.Query(context.Background(), &queryResult, variables)
		if err != nil {
			return nil, err
		}

		for _, node := range queryResult.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.Nodes {
			user := node.User
			if user.Id != nil {
				member := Member{
					Id:              user.Id.(string),
					Login:           string(user.Login),
					DatabaseId:      int64(user.DatabaseId),
					EnterpriseEmail: string(node.SamlIdentity.Username),
				}
				result.Members = append(result.Members, member)
			}
		}

		if !queryResult.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.HasNextPage {
			break
		}

		after = &queryResult.Enterprise.OwnerInfo.SamlIdentityProvider.ExternalIdentities.PageInfo.EndCursor
	}

	fmt.Printf("Time taken to fetch enterprise members : %v\n", time.Since(start))
	return &result, nil
}

func GetRepositoryProjects(owner string, name string, token string) (*GetRepositoryProjectsResult, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var result GetRepositoryProjectsResult
	var queryResult GetRepositoryProjectsQuery

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
	err := client.Query(context.Background(), &queryResult, variables)
	if err != nil {
		return nil, err
	}

	result.ProjectUrl = string(queryResult.Repository.ProjectsUrl)

	for _, project := range queryResult.Repository.ProjectsV2.Nodes {
		result.Projects = append(result.Projects, Project{
			Databaseid: int64(project.DatabaseId),
			Url:        string(project.Url),
			Title:      string(project.Title),
			CreatedAt:  project.CreatedAt.Time,
			UpdatedAt:  project.UpdatedAt.Time,
		})
	}

	return &result, nil
}

func GetUserByLogin(login string, token string) (*GetUserByLoginResult, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	httpClient.Transport = &customTransport{Transport: httpClient.Transport}

	client := githubv4.NewClient(httpClient)

	var result GetUserByLoginResult
	var queryResult GetUserByLoginQuery

	variables := map[string]interface{}{
		"login": githubv4.String(login),
	}
	err := client.Query(context.Background(), &queryResult, variables)
	if err != nil {
		return nil, err
	}

	if queryResult.User.ID == nil {
		return nil, fmt.Errorf("node ID of %s is not found", queryResult.User.Login)
	}

	result.User = User{
		Id:         queryResult.User.ID.(string),
		DatabaseId: int64(queryResult.User.DatabaseId),
		Login:      string(queryResult.User.Login),
	}

	return &result, nil
}

type GetOrganizationsByGitHubNameQuery struct {
	User struct {
		Organizations struct {
			Edges []struct {
				Node struct {
					Login      githubv4.String
					DatabaseId githubv4.Int
				}
			}
		} `graphql:"organizations(first: 90)"`
	} `graphql:"user(login: $login)"`
}

// Query structs
type GetOrganizationsWithinEnterpriseQuery struct {
	Enterprise struct {
		Organization struct {
			Nodes []struct {
				Login      githubv4.String
				DatabaseId githubv4.Int
			}
			PageInfo PageInfo
		} `graphql:"organizations(first: 100, after: $cursor)"`
	} `graphql:"enterprise(slug: $enterprise)"`
}

type GetMembersByEnterpriseQuery struct {
	Enterprise struct {
		OwnerInfo struct {
			SamlIdentityProvider struct {
				ExternalIdentities struct {
					Nodes []struct {
						SamlIdentity struct {
							Username githubv4.String
						}
						User struct {
							Id         githubv4.ID
							DatabaseId githubv4.Int
							Login      githubv4.String
						}
					}
					PageInfo PageInfo
				} `graphql:"externalIdentities(first: 100, after: $after)"`
			} `graphql:"samlIdentityProvider"`
		} `graphql:"ownerInfo"`
	} `graphql:"enterprise(slug: $enterprise)"`
}

type GetRepositoryProjectsQuery struct {
	Repository struct {
		ProjectsUrl githubv4.String
		ProjectsV2  struct {
			Nodes []struct {
				DatabaseId githubv4.Int
				Title      githubv4.String
				Id         githubv4.ID
				Url        githubv4.String
				CreatedAt  githubv4.DateTime
				UpdatedAt  githubv4.DateTime
			}
		} `graphql:"projectsV2(first: 100)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type GetUserByLoginQuery struct {
	User struct {
		ID         githubv4.ID
		DatabaseId githubv4.Int
		Login      githubv4.String
	} `graphql:"user(login: $login)"`
}

type PageInfo struct {
	EndCursor   githubv4.String
	HasNextPage bool
}

// Result structs

type GetOrganizationsByGithubNameResult struct {
	Organizations []Organization
}

type GetOrganizationsWithinEnterpriseResult struct {
	Organizations []Organization
}

type GetMembersByEnterpriseResult struct {
	Members []Member
}

type GetRepositoryProjectsResult struct {
	ProjectUrl string
	Projects   []Project
}

type GetUserByLoginResult struct {
	User
}

// Structs
type Member struct {
	Id              string
	Login           string
	DatabaseId      int64
	EnterpriseEmail string
}

type Organization struct {
	Login      string
	DatabaseId int64
}

type Project struct {
	Databaseid int64
	Url        string
	Title      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type User struct {
	Id         string
	DatabaseId int64
	Login      string
}
