package merger

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
)

// Organisation represents a GitHub Organisation
// While an org has teams I believe it is better
// to capture this data at repo level and then only
// target the teams that are actually in use as opposed
// to all teams arbitrarily
type Organisation struct {
	Name        string
	FullName    string
	Description string
	URL         string
	Email       string
	Members     []Member
}

func (h *Handler) orgDetails(ctx context.Context, name string) (Organisation, error) {
	var organisation Organisation

	// Connect to git
	org, _, err := h.clientRest.Organizations.Get(ctx, name)

	if err != nil {
		return organisation, err
	}
	organisation = Organisation{
		Name:        org.GetLogin(),
		FullName:    org.GetName(),
		Description: org.GetDescription(),
		URL:         org.GetURL(),
		Email:       org.GetEmail(),
	}

	orgMembers, err := h.orgMembers(ctx)
	if err != nil {
		return organisation, err
	}
	organisation.Members = orgMembers
	return organisation, nil
}

func (h *Handler) orgMembers(ctx context.Context) ([]Member, error) {
	h.log.Debugf("Gathering Org Members")
	opts := &github.ListMembersOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allMembers []Member
	for {
		opts.Page = page
		members, resp, err := h.clientRest.Organizations.ListMembers(ctx, h.config.SourceOrg.Name, opts)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			allMembers = append(allMembers, Member{
				Login: member.GetLogin(),
				Email: member.GetEmail(),
			})
		}

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return allMembers, nil
}

func (h *Handler) orgRepos(ctx context.Context) ([]Repository, error) {
	opts := github.RepositoryListByOrgOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allRepos []Repository
	for {
		opts.Page = page
		repos, resp, err := h.clientRest.Repositories.ListByOrg(ctx, h.config.SourceOrg.Name, &opts)
		if err != nil {
			return nil, err
		}
		for _, repo := range repos {
			h.log.Debugf("Gathering Repo Details: %s", repo.GetName())
			r := Repository{
				Name:        repo.GetName(),
				Description: repo.GetDescription(),
				URL:         repo.GetURL(),
				Private:     repo.GetPrivate(),
				PushedAt:    repo.GetPushedAt().String(),
			}
			t, err := h.repoTeams(ctx, repo.GetName())
			fmt.Println("Teams")
			fmt.Println(t)
			if err != nil {
				return nil, err
			}
			r.Teams = t

			c, err := h.repoCollaborators(ctx, repo.GetName())
			if err != nil {
				return nil, err
			}
			r.Collaborators = c

			con, err := h.repoContributors(ctx, repo.GetName())
			if err != nil {
				return nil, err
			}
			r.Contributors = con

			allRepos = append(allRepos, r)
		}

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return allRepos, nil
}
