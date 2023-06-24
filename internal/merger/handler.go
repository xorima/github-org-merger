package merger

import (
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"github.com/shurcooL/githubv4"
	"github.com/xorima/github-org-merger/internal/config"
	"github.com/youshy/logger"
	"go.uber.org/zap"
)

type OrganisationInformation struct {
	Organisation Organisation
	Repositories []Repository
	SeenTeams    []Team
}

type Member struct {
	Login string
	Email string
}

type Handler struct {
	config        *config.Config
	clientRest    *github.Client
	clientGraphql *githubv4.Client
	teamCache     map[string]Team
	log           *zap.SugaredLogger
	ctx           context.Context
}

func NewHandler(config *config.Config) *Handler {

	log := logger.NewLogger(logger.DEBUG, false)

	return &Handler{
		config:        config,
		clientRest:    NewGithubClientPAT(context.Background(), config.GithubToken),
		clientGraphql: NewGithubGraphqlClientPAT(context.Background(), config.GithubToken),
		teamCache:     make(map[string]Team),
		log:           log,
		ctx:           context.Background(),
	}

}

func (h *Handler) Gather() {
	h.log.Debugf("Running on Org: %s", h.config.SourceOrg.Name)
	var orgInfo OrganisationInformation
	h.log.Debugf("Gathering Org Details")
	org, err := h.orgDetails(h.ctx, h.config.SourceOrg.Name)
	if err != nil {
		panic(err)
	}
	orgInfo.Organisation = org
	h.log.Debugf("Gathering Repo Details")
	repos, err := h.orgRepos(h.ctx)
	if err != nil {
		panic(err)
	}
	orgInfo.Repositories = repos
	var teams []Team
	for _, v := range h.teamCache {
		teams = append(teams, v)
	}

	orgInfo.SeenTeams = teams
	h.printJson(orgInfo, orgInfo.Organisation.Name)
}

func (h *Handler) Plan() {
	if (h.config.SingleRepository == "" && !h.config.AllRepositories) || (h.config.SingleRepository != "" && h.config.AllRepositories) {
		panic("Must set either --repository or --all-repositories")
	}
	var repos []Repository
	if h.config.SingleRepository != "" {
		repos = h.planSingleRepo(h.config.SingleRepository)
	} else {
		repos = h.planAllRepos()
	}
	for _, r := range repos {
		h.log.Infof("For Repo %s gathering teams", r.Name)
		t, err := h.repoTeams(h.ctx, r.Name)
		if err != nil {
			h.log.Panicf("Failed to get repo teams for %s, error: %s", r.Name, err.Error())
			panic(err)
		}
		r.Teams = t
		h.log.Infof("Cached %d teams used for permissions on repo %s", len(t), r.Name)
		for _, t := range t {
			h.teamCache[t.Name] = t
		}
		bt, err := h.getBranchProtectionGroups(h.ctx, config.AppConfig.SourceOrg.Name, r.Name)
		if err != nil {
			h.log.Panicf("Failed to get branch protection groups for %s, error: %s", r.Name, err.Error())
			panic(err)
		}
		r.BypassTeams = bt
	}
	plan := h.generatePlan(repos)
	h.printJson(plan, "plan")

	fmt.Println(repos)
	// Generate a plan json file

	// TODO: Scan destination org for existing repos with same name and prefix names if needed on repos to be transferred
	// TODO: Scan desired destination org for existing teams and prefix names if needed on teams to be transferred
	// TODO: Scan destination org for members who are in teams to be transferred and not in the destination org

}

func (h *Handler) planSingleRepo(repoName string) []Repository {
	repo, err := h.repoDetails(h.ctx, repoName)
	if err != nil {
		h.log.Panicf("Failed to get repo details for %s, error: %s", repoName, err.Error())
		panic(err)
	}
	return []Repository{repo}
}

func (h *Handler) planAllRepos() []Repository {
	repos, err := h.orgRepos(h.ctx)
	if err != nil {
		h.log.Panicf("Failed to get all repositories for %s, error: %s", h.config.SourceOrg.Name, err.Error())

		panic(err)
	}
	return repos
}
