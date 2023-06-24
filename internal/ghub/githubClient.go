package ghub

import (
	"context"
	"github.com/google/go-github/v50/github"
	"github.com/shurcooL/githubv4"
	"github.com/xorima/github-org-merger/internal/config"
	"go.uber.org/zap"
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

func NewGithubGraphqlClientPAT(ctx context.Context, accessToken string) *githubv4.Client {
	httpClient := newOauthClientAccessToken(ctx, accessToken)
	return githubv4.NewClient(httpClient)
}

type Handler struct {
	clientRest    *github.Client
	clientGraphql *githubv4.Client
	config        *config.Config
	log           *zap.SugaredLogger
	ctx           context.Context
}

func NewHandler(ctx context.Context, accessToken string, cfg *config.Config, log *zap.SugaredLogger) *Handler {
	return &Handler{
		clientRest:    NewGithubClientPAT(ctx, accessToken),
		clientGraphql: NewGithubGraphqlClientPAT(ctx, accessToken),
		config:        cfg,
		log:           log,
		ctx:           ctx,
	}
}
