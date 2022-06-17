package githubAPI

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"main/models"
	"main/pkg/envvar"
	ghmgmt "main/pkg/ghmgmtdb"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v42/github"
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

func CreatePrivateGitHubRepository(data models.TypNewProjectReqBody) (*github.Repository, error) {
	client := createClient(os.Getenv("GH_TOKEN"))
	owner := envvar.GetEnvVar("GH_PROJECT_OWNER", "ava-innersource")
	repoRequest := &github.TemplateRepoRequest{
		Name:        &data.Name,
		Owner:       &owner,
		Description: &data.Description,
		Private:     github.Bool(true),
	}

	repo, _, err := client.Repositories.CreateFromTemplate(context.Background(), "avanade", "avanade-template", repoRequest)
	if err != nil {
		return nil, err
	}

	AddCollaborator(data)
	return repo, nil
}

func AddCollaborator(data models.TypNewProjectReqBody) (*github.Response, error) {
	client := createClient(os.Getenv("GH_TOKEN"))
	owner := envvar.GetEnvVar("GH_PROJECT_OWNER", "ava-innersource")
	opts := &github.RepositoryAddCollaboratorOptions{
		Permission: "admin",
	}

	GHUser := ghmgmt.Users_Get_GHUser(data.Coowner)

	_, resp, err := client.Repositories.AddCollaborator(context.Background(), owner, data.Name, GHUser, opts)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func GetRepository(repoName string) (*github.Repository, error) {
	client := createClient(os.Getenv("GH_TOKEN"))
	owner := envvar.GetEnvVar("GH_PROJECT_OWNER", "ava-innersource")
	repo, _, err := client.Repositories.Get(context.Background(), owner, repoName)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func Repo_IsExisting(repoName string) (bool, error) {
	_, err := GetRepository(repoName)
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			return false, nil
		}
		return false, err
	}

	return true, nil
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
			FullName:    repo.GetFullName(),
			Name:        repo.GetName(),
			Link:        repo.GetHTMLURL(),
			Org:         org,
			Description: repo.GetDescription(),
			Private:     repo.GetPrivate(),
			Created:     repo.GetCreatedAt(),
		}
		repoList = append(repoList, r)
	}

	return repoList, nil
}

func SetProjectVisibility(projectName string, visibility string) error {
	client := &http.Client{}
	urlPath := fmt.Sprintf("https://api.github.com/repos/%s/%s", envvar.GetEnvVar("GH_PROJECT_OWNER", "Avanade"), projectName)
	postBody, _ := json.Marshal(map[string]string{
		"visibility": visibility,
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

type Repo struct {
	FullName    string           `json:"repoFullName"`
	Name        string           `json:"repoName"`
	Link        string           `json:"repoLink"`
	Org         string           `json:"org"`
	Description string           `json:"description"`
	Private     bool             `json:"private"`
	Created     github.Timestamp `json:"created"`
	IsArchived  bool             `json:"archived"`
}
