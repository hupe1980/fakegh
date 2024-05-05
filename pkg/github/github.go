package github

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/v61/github"
	gogithub "github.com/google/go-github/v61/github"
	"github.com/hupe1980/fakegh/pkg/util"
)

type Client struct {
	client *gogithub.Client
}

func New(token string) *Client {
	client := gogithub.NewClient(nil).WithAuthToken(token)

	return &Client{
		client: client,
	}
}

func (gh *Client) GetEmailByUsername(ctx context.Context, username string) (*string, error) {
	user, _, err := gh.client.Users.Get(ctx, username)
	if err != nil {
		return nil, err
	}

	if user.Email != nil {
		return user.Email, nil
	}

	repos, _, err := gh.client.Repositories.ListByUser(context.Background(), username, &github.RepositoryListByUserOptions{
		Sort:      "pushed",
		Direction: "desc",
	})
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		commits, _, err := gh.client.Repositories.ListCommits(context.Background(), username, *repo.Name, &github.CommitsListOptions{
			Author: username,
		})
		if err != nil {
			return nil, err
		}

		for _, commit := range commits {
			if commit.Commit != nil {
				if commit.Commit.Committer != nil {
					if commit.Commit.Committer.Email != nil {
						if *commit.Commit.Committer.Email != "noreply@github.com" {
							return commit.Commit.Committer.Email, nil
						}
					}
				}
			}
		}
	}

	return nil, fmt.Errorf("cannot find mail for user %s", username)
}

func (gh *Client) OpenIssue(ctx context.Context, url string) (*int, error) {
	owner, repo := extractOwnerAndRepo(url)

	issue, _, err := gh.client.Issues.Create(ctx, owner, repo, &gogithub.IssueRequest{
		Title: util.PTR("Test"),
		Body:  util.PTR("Test"),
	})
	if err != nil {
		return nil, err
	}

	return issue.Number, nil
}

func (gh *Client) CloseIssue(ctx context.Context, url string, number int) error {
	owner, repo := extractOwnerAndRepo(url)

	_, _, err := gh.client.Issues.Edit(ctx, owner, repo, number, &gogithub.IssueRequest{
		State: util.PTR("closed"),
	})
	if err != nil {
		return err
	}

	return nil
}

func extractOwnerAndRepo(url string) (string, string) {
	re := regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+)`)
	match := re.FindStringSubmatch(url)
	if len(match) < 3 {
		return "", ""
	}

	return match[1], match[2]
}
