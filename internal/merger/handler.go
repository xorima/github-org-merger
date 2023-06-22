package merger

import (
	"context"
	"encoding/json"
	"fmt"
	"github-org-merger/internal/config"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/youshy/logger"
	"go.uber.org/zap"
)

type OrganisationInformation struct {
	Organisation Organisation
	Repositories []Repository
	SeenTeams    []Team
}

// Organisation represents a github Organisation
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
type Member struct {
	Login string
	Email string
}
type Team struct {
	Name        string
	Description string
	URL         string
	Parent      string
}
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

type Handler struct {
	config    *config.Config
	client    *github.Client
	teamCache map[string]Team
	log       *zap.SugaredLogger
}

func NewHandler(config *config.Config) *Handler {

	log := logger.NewLogger(logger.DEBUG, false)

	return &Handler{
		config:    config,
		client:    NewGithubClientPAT(context.Background(), config.GithubToken),
		teamCache: make(map[string]Team),
		log:       log,
	}

}

func (h *Handler) Handle() {
	h.log.Debugf("Running on Org: %s", h.config.SourceOrg.Name)
	var orgInfo OrganisationInformation
	h.log.Debugf("Gathering Org Details")
	org, err := h.orgDetails()
	if err != nil {
		panic(err)
	}
	orgInfo.Organisation = org
	h.log.Debugf("Gathering Repo Details")
	repos, err := h.orgRepos()
	if err != nil {
		panic(err)
	}
	orgInfo.Repositories = repos

	// TODO: Add teams
	var teams []Team
	for _, v := range h.teamCache {
		teams = append(teams, v)
	}

	orgInfo.SeenTeams = teams
	h.printJson(orgInfo)
}

func (h *Handler) printJson(orgInfo OrganisationInformation) {
	h.log.Debugf("Printing JSON to screen")
	// convert to json
	j, err := json.Marshal(orgInfo)
	if err != nil {
		panic(err)
	}
	// save to disk using org name
	h.log.Debugf("Saving JSON to disk")
	fmt.Println(orgInfo.Organisation.Name)
	err = h.saveJson(j, orgInfo.Organisation.Name)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) saveJson(json []byte, orgName string) error {
	h.log.Debugf("Saving to file as %s", orgName)
	return os.WriteFile(fmt.Sprintf("%s.json", orgName), json, 0644)
}

func (h *Handler) githubListOptsDefaults() github.ListOptions {
	return github.ListOptions{PerPage: 100}
}

func (h *Handler) orgDetails() (Organisation, error) {
	var organisation Organisation

	// Connect to git
	org, _, err := h.client.Organizations.Get(context.Background(), h.config.SourceOrg.Name)

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

	orgMembers, err := h.orgMembers()
	if err != nil {
		return organisation, err
	}
	organisation.Members = orgMembers
	return organisation, nil
}

func (h *Handler) orgMembers() ([]Member, error) {
	h.log.Debugf("Gathering Org Members")
	opts := &github.ListMembersOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allMembers []Member
	for {
		opts.Page = page
		members, resp, err := h.client.Organizations.ListMembers(context.Background(), h.config.SourceOrg.Name, opts)
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

func (h *Handler) teamDetails() ([]Team, error) {
	opts := h.githubListOptsDefaults()
	page := 1
	var allTeams []Team
	for {
		opts.Page = page
		teams, resp, err := h.client.Teams.ListTeams(context.Background(), h.config.SourceOrg.Name, &opts)
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

func (h *Handler) orgRepos() ([]Repository, error) {
	opts := github.RepositoryListByOrgOptions{
		ListOptions: h.githubListOptsDefaults(),
	}
	page := 1
	var allRepos []Repository
	for {
		opts.Page = page
		repos, resp, err := h.client.Repositories.ListByOrg(context.Background(), h.config.SourceOrg.Name, &opts)
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
			t, err := h.repoTeams(repo.GetName())
			if err != nil {
				return nil, err
			}
			r.Teams = t

			c, err := h.repoCollaborators(repo.GetName())
			if err != nil {
				return nil, err
			}
			r.Collaborators = c

			con, err := h.repoContributors(repo.GetName())
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
