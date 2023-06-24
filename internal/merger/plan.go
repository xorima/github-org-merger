package merger

import "fmt"

type Plan struct {
	SourceOrganisation      Organisation
	DestinationOrganisation Organisation
	Create                  PlanCreate
	Migrate                 PlanMigrate
	Update                  PlanUpdate
}

type PlanCreate struct {
	Teams PlanCreateTeams
}

type PlanCreateTeams struct {
	DesiredTeams []Team
	CurrentTeams []Team
}

type PlanMigrate struct {
	DesiredRepositories []Repository
}

type PlanUpdate struct {
	ProtectionBypasses PlanUpdateProtectionBypasses
}

type PlanUpdateProtectionBypasses struct {
	Repositories []PlanUpdateProtectionRepositories
}
type PlanUpdateProtectionRepositories struct {
	Repository        Repository
	DesiredProtection ProtectionBypasses
	CurrentProtection ProtectionBypasses
}

func (h *Handler) generatePlan(repos []Repository) Plan {
	sourceOrg, err := h.orgDetails(h.ctx, h.config.SourceOrg.Name)
	if err != nil {
		h.log.Panicf("Failed to get org details for %s, error: %s", h.config.SourceOrg.Name, err.Error())
		panic(err)
	}

	destinationOrg, err := h.orgDetails(h.ctx, h.config.DestinationOrg.Name)
	if err != nil {
		h.log.Panicf("Failed to get org details for %s, error: %s", h.config.DestinationOrg.Name, err.Error())
		panic(err)
	}

	createTeams := PlanCreateTeams{}
	fmt.Println(h.teamCache)
	// Foreach team in cache add to plan
	for _, repo := range repos {
		for _, t := range repo.Teams {
			createTeams.CurrentTeams = append(createTeams.CurrentTeams, t)
			tmp := Team{
				Name:        h.getTeamName(t.Name),
				Description: t.Description,
				URL:         t.URL,
				Parent:      t.Parent,
				Members:     t.Members,
			}
			createTeams.DesiredTeams = append(createTeams.DesiredTeams, tmp)
		}
	}

	pc := PlanCreate{
		Teams: createTeams,
	}
	pm := PlanMigrate{
		DesiredRepositories: repos,
	}

	pu := PlanUpdate{
		ProtectionBypasses: PlanUpdateProtectionBypasses{},
	}

	for _, r := range repos {
		for _, t := range r.BypassTeams {
			tmp := PlanUpdateProtectionRepositories{
				Repository: r,
			}
			tmp.CurrentProtection = t

			tmpBypassProtect := ProtectionBypasses{
				BranchPattern: t.BranchPattern,
			}
			for _, prt := range t.BypassPullRequestTeams {
				tmpBypassProtect.BypassPullRequestTeams = append(tmpBypassProtect.BypassPullRequestTeams, h.getTeamName(prt))
			}
			for _, fp := range t.BypassForcePush {
				tmpBypassProtect.BypassForcePush = append(tmpBypassProtect.BypassForcePush, h.getTeamName(fp))
			}

			pu.ProtectionBypasses.Repositories = append(pu.ProtectionBypasses.Repositories, tmp)
		}
	}

	plan := Plan{
		SourceOrganisation:      sourceOrg,
		DestinationOrganisation: destinationOrg,
		Create:                  pc,
		Migrate:                 pm,
		Update:                  pu,
	}
	return plan
}

// Checks the destination org for teams  and ensures the current proposed name
// does not clash
func (h *Handler) getTeamName(teamName string) string {
	return teamName
}
