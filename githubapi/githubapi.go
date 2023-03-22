// githubapi/githubapi.go
package githubapi

import (
	"context"

	"github.com/google/go-github/v50/github"
	"github.com/mhassaninmsft/gogitcrawler/models"
	"golang.org/x/oauth2"
)

const perPage = 100

type Client struct {
	client *github.Client
}

func NewClient(token string) *Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return &Client{client: github.NewClient(tc)}
}

func (c *Client) GetContributors(repoOwner, repoName string) ([]*models.Contributor, error) {
	ctx := context.Background()
	var contributors []*models.Contributor
	opts := &github.ListContributorsOptions{ListOptions: github.ListOptions{PerPage: perPage}}
	// opts := &github.ListContributorsOptions{PerPage: perPage}
	for {
		contribs, resp, err := c.client.Repositories.ListContributors(ctx, repoOwner, repoName, opts)
		if err != nil {
			return nil, err
		}
		for _, contrib := range contribs {
			contributors = append(contributors, &models.Contributor{
				ID:    contrib.GetID(),
				Login: contrib.GetLogin(),
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return contributors, nil
}

func (c *Client) GetPublicRepos(username string) ([]*models.Repo, error) {
	ctx := context.Background()
	var repos []*models.Repo
	opts := &github.RepositoryListOptions{Type: "public", ListOptions: github.ListOptions{PerPage: perPage}	}
	for {
		reposList, resp, err := c.client.Repositories.List(ctx, username, opts)
		if err != nil {
			return nil, err
		}
		for _, repo := range reposList {
			repos = append(repos, &models.Repo{
				ID:       repo.GetID(),
				Name:     repo.GetName(),
				FullName: repo.GetFullName(),
			})
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return repos, nil
}
