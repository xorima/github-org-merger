package merger

import (
	"context"
	"github.com/xorima/pointerhelpers"
)

type Team struct {
	Name        string
	Description string
	URL         string
	Parent      string
	Members     []Member
	Permission  *string
}

func (t *Team) GetPermission() (bool, string) {
	if t.Permission == nil {
		return false, ""
	}
	return true, pointerhelpers.StringValue(t.Permission)
}

func (h *Handler) teamDetails(ctx context.Context, orgName string) ([]Team, error) {
	opts := h.githubListOptsDefaults()
	page := 1
	var allTeams []Team
	for {
		opts.Page = page
		teams, resp, err := h.clientRest.Teams.ListTeams(ctx, orgName, &opts)
		if err != nil {
			return nil, err
		}
		for _, team := range teams {
			t := Team{
				Name:        team.GetName(),
				Description: team.GetDescription(),
				URL:         team.GetURL(),
				Parent:      team.GetParent().GetName(),
				Permission:  team.Permission,
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
