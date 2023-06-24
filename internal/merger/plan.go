package merger

import (
	"fmt"
	"github.com/xorima/github-org-merger/internal/ghub"
)

type TeamMigration struct {
	Create []ghub.Team
	Reuse  []ghub.Team
	Issues []string
}

type RepoMigration struct {
	Migrate []ghub.Repository
	Issues  []string
}

func (h *Handler) teamsToMigrate(src, dest ghub.Organisation) (TeamMigration, map[string]string, []ghub.Member) {

	allSlugs := make(map[string]bool)

	for _, repo := range src.Repositories {
		h.log.Infof("Checking access teams for %s", repo.Name)
		for _, accessTeam := range repo.AccessTeams {
			allSlugs[accessTeam.Slug] = true
		}
		h.log.Infof("Checking bypass teams for %s", repo.Name)
		for _, bypassTeams := range repo.ProtectedBranchBypasses {
			for _, bypassTeamPR := range bypassTeams.BypassPullRequestTeams {
				allSlugs[bypassTeamPR] = true
			}
			for _, bypassTeamFP := range bypassTeams.BypassForcePush {
				allSlugs[bypassTeamFP] = true
			}
		}
	}
	missingMembers := h.teamMembersMissingInDestOrg(allSlugs, src, dest)
	teamsMigration, teamRenames := h.teamsMigration(allSlugs, src, dest, missingMembers)
	return teamsMigration, teamRenames, missingMembers
}

func (h *Handler) teamMembersMissingInDestOrg(slugs map[string]bool, src, dest ghub.Organisation) []ghub.Member {
	membersSeenSrc := make(map[string]ghub.Member)
	for slug, _ := range slugs {
		for _, member := range src.Teams[slug].Members {
			fmt.Printf("**** SEEN MEMBER **** %s\n", member.Login)
			membersSeenSrc[member.Login] = member
		}
	}

	var allDestMembersToday = make(map[string]ghub.Member)
	for _, member := range dest.Members {
		fmt.Printf("**** SEEN DEST MEMBER **** %s\n", member.Login)

		allDestMembersToday[member.Login] = member
	}

	var missingMembers []ghub.Member
	for login, _ := range membersSeenSrc {
		fmt.Println("checking if member is in dest org", login)
		if _, ok := allDestMembersToday[login]; !ok {
			fmt.Println("member is missing from dest org", login)
			missingMembers = append(missingMembers, membersSeenSrc[login])
			continue
		}
		fmt.Println("member is in dest org", login)
	}

	return missingMembers
}

func (h *Handler) teamsMigration(slugs map[string]bool, src, dest ghub.Organisation, missingMembers []ghub.Member) (TeamMigration, map[string]string) {
	var teamMigration = TeamMigration{}
	var teamRenames = make(map[string]string)
	for slug, _ := range slugs {
		fmt.Println("SLUG", slug)
		// if team does not exist in destination, create it
		if _, ok := dest.Teams[slug]; !ok {
			teamMigration.Create = append(teamMigration.Create, src.Teams[slug])
			continue
		}
		// if team does exist, check its members
		for _, member := range src.Teams[slug].Members {
			for _, missMember := range missingMembers {
				if missMember.Login == member.Login {
					msg := fmt.Sprintf(" Member %s is not a member in the destination org, but they are a part of %s which is an in use team on the source org. Ensure they are added to the destination org OR removed from the source team %s", member.Login, slug, slug)
					h.log.Infof(msg)
					teamMigration.Issues = append(teamMigration.Issues, msg)
					continue
				}

			}
			// does the destination team contain this member? If not we will need to create a prefixed team
			destTeam := dest.Teams[slug]
			for _, destMember := range destTeam.Members {
				if destMember.Login == member.Login {

				}
			}
			if _, ok := dest.Teams[slug].Members[member.Login]; !ok {
				h.log.Infof("Team %s does not contain member %s, so we will create a new team with org-prefix called: %s, no action required on your part.", slug, member.Login, src.Name)
				create := src.Teams[slug]
				create.Name = h.getPrefixedName(create.Name, src)
				teamMigration.Create = append(teamMigration.Create, src.Teams[slug])
				teamRenames[slug] = create.Name
				continue
			}
		}
	}
	return teamMigration, teamRenames
}

func (h *Handler) getPrefixedName(name string, org ghub.Organisation) string {
	return fmt.Sprintf("%s-%s", org.Name, name)
}

func (h *Handler) ValidateReposToMigrate(src, dest ghub.Organisation) RepoMigration {

	var response RepoMigration
	for _, repo := range src.Repositories {
		if _, ok := dest.Repositories[repo.Name]; ok {
			response.Issues = append(response.Issues, fmt.Sprintf("Repository %s exists on the destination org, please rename to a unique name", repo.Name))
			continue
		}
		response.Migrate = append(response.Migrate, repo)
	}
	return response
}
