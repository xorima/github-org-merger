package ghub

import (
	"context"
	"github.com/google/go-github/v50/github"
	"github.com/xorima/pointerhelpers"
)

type Member struct {
	Login string
	Email string
}
type Team struct {
	Name        string
	Description string
	URL         string
	Parent      string
	Members     map[string]Member
	Permission  *string
}

func (t *Team) GetPermission() (bool, string) {
	if t.Permission == nil {
		return false, ""
	}
	return true, pointerhelpers.StringValue(t.Permission)
}

func teamName(team *github.Team) string {
	return team.GetSlug()
}

func (h *Handler) teamDetails(ctx context.Context, orgName string) (map[string]Team, error) {
	opts := h.githubListOptsDefaults()
	page := 1
	var allTeams = make(map[string]Team)

	for {
		opts.Page = page
		teams, resp, err := h.clientRest.Teams.ListTeams(ctx, orgName, &opts)
		if err != nil {
			return nil, err
		}
		for _, team := range teams {
			members, err := h.teamMembers(ctx, orgName, team.GetSlug())
			if err != nil {
				return nil, err
			}
			t := Team{
				Name:        team.GetName(),
				Description: team.GetDescription(),
				URL:         team.GetURL(),
				Parent:      team.GetParent().GetName(),
				Permission:  team.Permission,
				Members:     members,
			}
			allTeams[teamName(team)] = t
		}
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}

	return allTeams, nil
}

func (h *Handler) teamMembers(ctx context.Context, org, slug string) (map[string]Member, error) {
	opts := github.TeamListTeamMembersOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allMembers = make(map[string]Member)

	for {
		opts.Page = page
		members, resp, err := h.clientRest.Teams.ListTeamMembersBySlug(ctx, org, slug, &opts)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			allMembers[member.GetLogin()] = Member{
				Login: member.GetLogin(),
				Email: member.GetEmail(),
			}
		}
		if resp.NextPage == 0 {
			break
		}
		page = resp.NextPage
	}
	return allMembers, nil
}
