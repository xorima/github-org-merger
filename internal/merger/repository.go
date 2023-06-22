package merger

import (
	"context"
	"github.com/google/go-github/v50/github"
	"strings"
)

type Repository struct {
	Name          string
	Description   string
	URL           string
	Private       bool
	Teams         []Team
	Collaborators []Member
	Contributors  []Member
	PushedAt      string
}

func (h *Handler) repoTeams(repo string) ([]Team, error) {
	h.log.Debugf("Gathering Repo Teams: %s", repo)
	opts := h.githubListOptsDefaults()
	page := 1
	var allTeams []Team
	for {
		opts.Page = page
		teams, resp, err := h.client.Repositories.ListTeams(context.Background(), h.config.SourceOrg.Name, repo, &opts)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				h.log.Warnf("No access to TEAMS for %s, this data will not be captured", repo)
				return nil, nil
			}
			return nil, err
		}
		for _, team := range teams {
			h.log.Debugf("Gathering Repo Team Details: %s", team.GetName())
			t := Team{
				Name:        team.GetName(),
				Description: team.GetDescription(),
				URL:         team.GetURL(),
				Parent:      team.GetParent().GetName(),
			}
			h.teamCache[t.Name] = t
			allTeams = append(allTeams, t)
		}
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

	return allTeams, nil

}

func (h *Handler) repoCollaborators(repo string) ([]Member, error) {
	h.log.Debugf("Gathering Repo Collaborators: %s", repo)
	opts := &github.ListCollaboratorsOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allCollaborators []Member
	for {
		opts.Page = page
		collaborators, resp, err := h.client.Repositories.ListCollaborators(context.Background(), h.config.SourceOrg.Name, repo, opts)
		if err != nil {
			if (strings.Contains(err.Error(), "404")) || (strings.Contains(err.Error(), "403")) {
				h.log.Warnf("No access to COLLABORATORS for %s, this data will not be captured", repo)
				return nil, nil
			}
			return nil, err
		}
		for _, collaborator := range collaborators {
			h.log.Debugf("Gathering Repo Collaborator Details: %s", collaborator.GetLogin())
			allCollaborators = append(allCollaborators, Member{
				Login: collaborator.GetLogin(),
				Email: collaborator.GetEmail(),
			})
		}
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return allCollaborators, nil
}

func (h *Handler) repoContributors(repo string) ([]Member, error) {
	h.log.Debugf("Gathering Repo Contributors: %s", repo)
	opts := &github.ListContributorsOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allContributors []Member
	for {
		opts.Page = page
		contributors, resp, err := h.client.Repositories.ListContributors(context.Background(), h.config.SourceOrg.Name, repo, opts)
		if err != nil {
			return nil, err
		}
		for _, contributor := range contributors {
			h.log.Debugf("Gathering Repo Contributor Details: %s", contributor.GetLogin())
			allContributors = append(allContributors, Member{
				Login: contributor.GetLogin(),
				Email: contributor.GetEmail(),
			})
		}
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return allContributors, nil
}
