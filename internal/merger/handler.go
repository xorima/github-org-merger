package merger

import (
	"context"
	"fmt"
	"github-org-merger/internal/config"
	"github.com/google/go-github/v50/github"
)

type Handler struct {
	config *config.Config
	client *github.Client
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{
		config: config,
		client: NewGithubClientPAT(context.Background(), config.GithubToken),
	}
}

func (h *Handler) Handle() {
	fmt.Printf("Running on Org: %s\n", config.AppConfig.SourceOrg.Name)
	h.orgDetails()
}

func (h *Handler) orgDetails() error {
	// Connect to git
	org, _, err := h.client.Organizations.Get(context.Background(), h.config.SourceOrg.Name)
	if err != nil {
		return err
	}
	fmt.Println(org.GetName())
	fmt.Println(org.GetDescription())
	fmt.Println(org.GetURL())
	fmt.Println(org.GetEmail())

	h.repoDetails()
	return h.teamDetails()
}

func (h *Handler) teamDetails() error {
	teams, _, err := h.client.Teams.ListTeams(context.Background(), h.config.SourceOrg.Name, nil)
	if err != nil {
		return err
	}
	for _, team := range teams {
		fmt.Println(team.GetName())
		fmt.Println(team.GetDescription())
		fmt.Println(team.GetURL())
		fmt.Println(team.GetPermissions())
		fmt.Println(team.GetParent())
	}
	return nil
}

func (h *Handler) repoDetails() error {
	repos, _, err := h.client.Repositories.ListByOrg(context.Background(), h.config.SourceOrg.Name, nil)
	if err != nil {
		return err
	}
	for _, repo := range repos {
		fmt.Printf("Repo: %s\n", repo.GetName())
		fmt.Println(repo.GetName())
		fmt.Println(repo.GetDescription())
		fmt.Println(repo.GetURL())
		fmt.Println(repo.GetDefaultBranch())
		fmt.Println(repo.GetArchived())
		fmt.Println(repo.GetPrivate())
		fmt.Println(repo.GetPermissions())
		fmt.Println(repo.GetOwner())
		fmt.Println(repo.GetPushedAt())
		repo.GetPermissions()
		contribs, _, err := h.client.Repositories.ListContributors(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), nil)
		if err != nil {
			return err
		}
		for _, contrib := range contribs {
			fmt.Println("--------------")
			fmt.Println(contrib.GetLogin())
			fmt.Println(contrib.GetID())
			fmt.Println(contrib.GetURL())
			fmt.Println(contrib.GetAvatarURL())
			fmt.Println(contrib.GetType())
			fmt.Println("--------------")
		}

		collabs, _, err := h.client.Repositories.ListCollaborators(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), nil)
		if err != nil {
			return err
		}
		for _, collab := range collabs {
			fmt.Println("***********")
			fmt.Println(collab.GetLogin())
			fmt.Println(collab.GetID())
			fmt.Println(collab.GetURL())
			fmt.Println(collab.GetAvatarURL())
			fmt.Println(collab.GetType())
			fmt.Println("***********")
		}

		teams, _, err := h.client.Repositories.ListTeams(context.Background(), repo.GetOwner().GetLogin(), repo.GetName(), nil)
		if err != nil {
			return err
		}
		for _, team := range teams {
			fmt.Println("^^^^^^^^^^^^")
			fmt.Println(team.GetName())
			fmt.Println(team.GetDescription())
			fmt.Println(team.GetURL())
			fmt.Println(team.GetPermissions())
			fmt.Println(team.GetParent())
			fmt.Println("^^^^^^^^^^^^")

		}

		fmt.Println("#####")
	}
	return nil
}
