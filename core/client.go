package core

import (
	"context"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func newGithubClient(token string) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)
	return client
}
