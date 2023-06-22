package merger

import "github.com/google/go-github/v50/github"

func (h *Handler) githubListOptsDefaults() github.ListOptions {
	return github.ListOptions{PerPage: 100}
}
