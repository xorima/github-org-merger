package ghub

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
	Name         string
	FullName     string
	Description  string
	URL          string
	Email        string
	Members      []Member
	Teams        map[string]Team
	Repositories map[string]Repository
}

func (h *Handler) OrgDetails(ctx context.Context, name string, isSourceOrg bool) (Organisation, error) {
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

	organisation.Members, err = h.orgMembers(ctx, name)
	if err != nil {
		return organisation, err
	}
	organisation.Teams, err = h.teamDetails(ctx, name)
	if err != nil {
		return organisation, err
	}
	organisation.Repositories, err = h.orgRepos(ctx, name, isSourceOrg)
	if err != nil {
		return organisation, err
	}

	return organisation, nil
}

func (h *Handler) orgMembers(ctx context.Context, orgName string) ([]Member, error) {
	h.log.Debugf("Gathering Org Members")
	opts := &github.ListMembersOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allMembers []Member
	for {
		opts.Page = page
		members, resp, err := h.clientRest.Organizations.ListMembers(ctx, orgName, opts)
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

func (h *Handler) orgRepos(ctx context.Context, orgName string, includeTeamDetails bool) (map[string]Repository, error) {
	opts := github.RepositoryListByOrgOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allRepos = make(map[string]Repository)
	for {
		opts.Page = page
		repos, resp, err := h.clientRest.Repositories.ListByOrg(ctx, orgName, &opts)
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
			if !includeTeamDetails {
				allRepos[repo.GetName()] = r
				continue
			}

			t, err := h.repoTeams(ctx, repo.GetName())
			fmt.Println("Teams")
			fmt.Println(t)
			if err != nil {
				return nil, err
			}
			r.AccessTeams = t

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

			protect, err := h.getBranchProtectionGroups(ctx, orgName, repo.GetName())
			if err != nil {
				return nil, err
			}
			r.ProtectedBranchBypasses = protect

			allRepos[repo.GetName()] = r
		}

		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return allRepos, nil
}
