package githubAPI

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
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

func GetMembersByEnterprise(enterprise string, token string) (*GetMembersByEnterpriseResult, error) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

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
			member := Member{
				Login:           string(user.Login),
				DatabaseId:      int64(user.DatabaseId),
				EnterpriseEmail: string(node.SamlIdentity.Username),
			}

			result.Members = append(result.Members, member)
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

type PageInfo struct {
	EndCursor   githubv4.String
	HasNextPage bool
}

// Result structs
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

// Structs
type Member struct {
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
