package merger

import (
	"context"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
	"net/http"
)

func NewGithubClientPAT(ctx context.Context, accessToken string) *github.Client {
	httpClient := newOauthClientAccessToken(ctx, accessToken)
	return github.NewClient(httpClient)
}

func newOauthClientAccessToken(ctx context.Context, accessToken string) *http.Client {
	c := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	return oauth2.NewClient(ctx, c)
}
