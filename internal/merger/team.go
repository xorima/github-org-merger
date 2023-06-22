package merger

import "context"

type Team struct {
	Name        string
	Description string
	URL         string
	Parent      string
}

func (h *Handler) teamDetails(ctx context.Context) ([]Team, error) {
	opts := h.githubListOptsDefaults()
	page := 1
	var allTeams []Team
	for {
		opts.Page = page
		teams, resp, err := h.client.Teams.ListTeams(ctx, h.config.SourceOrg.Name, &opts)
		if err != nil {
			return nil, err
		}
		for _, team := range teams {
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
