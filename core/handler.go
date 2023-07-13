package core

import (
	"github.com/shurcooL/githubv4"
)

type Handler struct {
	client *githubv4.Client
}

func NewHandler(token string) *Handler {
	return &Handler{
		client: newGithubClient(token),
	}
}
