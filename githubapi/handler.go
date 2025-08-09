package githubapi

import (
	"context"
	"net/http"

	"github.com/google/go-github/v55/github"
)

type GithubClient struct {
	BaseURL string
	Client  *github.Client
	token   string
}

// tokenTransport injects an Authorization header with a static token
type tokenTransport struct {
	Token string
	Base  http.RoundTripper
}

func (t *tokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	clone := req.Clone(req.Context())
	if t.Token != "" {
		clone.Header.Set("Authorization", "token "+t.Token)
	}
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(clone)
}

func NewClient(token string) *GithubClient {
	httpClient := &http.Client{Transport: &tokenTransport{Token: token}}
	return &GithubClient{
		BaseURL: "https://api.github.com",
		Client:  github.NewClient(httpClient),
		token:   token,
	}
}

func (c *GithubClient) ListRepoContents(owner, repo, path string) []*github.RepositoryContent {
	file, dir, _, err := c.Client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return nil
	}
	if file != nil {

		return []*github.RepositoryContent{file}
	}
	return dir
}

func (c *GithubClient) GetFile(owner, repo, path string) (*github.RepositoryContent, error) {
	file, _, _, err := c.Client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return nil, err
	}
	return file, nil
}
