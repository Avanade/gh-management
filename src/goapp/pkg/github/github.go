package githubAPI

import (
	"context"
	"main/models"
	"os"

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
	repoRequest := &github.TemplateRepoRequest{
		Name:        &data.Name,
		Owner:       &data.Coowner,
		Description: &data.Description,
		Private:     github.Bool(true),
	}

	repo, _, err := GitHubClient.Repositories.CreateFromTemplate(context.Background(), "avanade", "avanade-template", repoRequest)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
