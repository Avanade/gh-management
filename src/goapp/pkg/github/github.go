package githubAPI

import (
	"context"
	"main/models"
	"main/pkg/envvar"
	"os"
	"strings"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

var (
	GitHubClient *github.Client
)

func CreateClient() {
	// create github oauth client from token
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GH_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	GitHubClient = github.NewClient(tc)
}

func CreatePrivateGitHubRepository(data models.TypNewProjectReqBody) (*github.Repository, error) {
	owner := envvar.GetEnvVar("GH_PROJECT_OWNER", "ava-innersource")
	repoRequest := &github.TemplateRepoRequest{
		Name:        &data.Name,
		Owner:       &owner,
		Description: &data.Description,
		Private:     github.Bool(true),
	}

	repo, _, err := GitHubClient.Repositories.CreateFromTemplate(context.Background(), "avanade", "avanade-template", repoRequest)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func GetRepository(repoName string) (*github.Repository, error) {
	owner := envvar.GetEnvVar("GH_PROJECT_OWNER", "ava-innersource")
	repo, _, err := GitHubClient.Repositories.Get(context.Background(), owner, repoName)
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
