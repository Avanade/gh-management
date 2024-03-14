package githubAPI

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

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

func GetRepositoryReadmeById(repoName string, visibility string) (string, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))

	if visibility == "" {
		return "", fmt.Errorf("invalid visibility")
	}

	var owner string

	if visibility != "Public" {
		owner = os.Getenv("GH_ORG_INNERSOURCE")
	} else {
		owner = os.Getenv("GH_ORG_OPENSOURCE")
	}

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
	organizations := []string{os.Getenv("GH_ORG_INNERSOURCE"), os.Getenv("GH_ORG_OPENSOURCE")}

	for _, org := range organizations {
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

func TransferRepository(repo string, owner string, newOwner string) (*github.Repository, error) {
	client := CreateClient(os.Getenv("GH_TOKEN"))
	opt := github.TransferRequest{NewOwner: newOwner}

	resp, _, err := client.Repositories.Transfer(context.Background(), owner, repo, opt)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func IsOrganizationMember(token, org, ghUser string) (bool, error) {
	client := CreateClient(token)
	isOrgMember, _, err := client.Organizations.IsMember(context.Background(), org, ghUser)
	return isOrgMember, err
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

func OrgListMembers(token string, org string) []*github.User {
	client := CreateClient(token)

	opts := &github.ListMembersOptions{ListOptions: github.ListOptions{PerPage: 30}}

	var members []*github.User

	for {
		listMembers, resp, err := client.Organizations.ListMembers(context.Background(), org, opts)
		if err != nil {
			log.Printf("ERROR : %s", err.Error())
			return nil
		}

		members = append(members, listMembers...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return members
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
