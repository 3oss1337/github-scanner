package githubapi

import (
	"context"

	"github.com/google/go-github/v55/github"
)

type GithubClient struct {
	BaseURL string
	Client  *github.Client
	token   string
}

func NewClient(token string) *GithubClient {
	return &GithubClient{
		BaseURL: "https://api.github.com",
		Client:  github.NewClient(nil),
		token:   token,
	}
}

func (c *GithubClient) GetRepoFiles(owner, repo, path string) []*github.RepositoryContent {
	files, _, err := c.Client.Repositories.ListContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return nil
	}
	return files
}

func (c *GithubClient) GetRepoContents(owner, repo, path string) []*github.RepositoryContent {
	contents, _, err := c.Client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return nil
	}
	return contents
}
