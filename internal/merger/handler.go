package merger

import (
	"context"
	"github.com/google/go-github/v50/github"
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
	config    *config.Config
	client    *github.Client
	teamCache map[string]Team
	log       *zap.SugaredLogger
	ctx       context.Context
}

func NewHandler(config *config.Config) *Handler {

	log := logger.NewLogger(logger.DEBUG, false)

	return &Handler{
		config:    config,
		client:    NewGithubClientPAT(context.Background(), config.GithubToken),
		teamCache: make(map[string]Team),
		log:       log,
		ctx:       context.Background(),
	}

}

func (h *Handler) Handle() {
	h.log.Debugf("Running on Org: %s", h.config.SourceOrg.Name)
	var orgInfo OrganisationInformation
	h.log.Debugf("Gathering Org Details")
	org, err := h.orgDetails(h.ctx)
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

	// TODO: Add teams
	var teams []Team
	for _, v := range h.teamCache {
		teams = append(teams, v)
	}

	orgInfo.SeenTeams = teams
	h.printJson(orgInfo)
}
