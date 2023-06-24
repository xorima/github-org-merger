package merger

import (
	"context"
	"fmt"
	"github.com/xorima/github-org-merger/internal/config"
	"github.com/xorima/github-org-merger/internal/ghub"
	"github.com/youshy/logger"
	"go.uber.org/zap"
)

type Handler struct {
	config *config.Config
	log    *zap.SugaredLogger
	ctx    context.Context
}

func NewHandler(config *config.Config) *Handler {

	log := logger.NewLogger(logger.DEBUG, false)

	return &Handler{
		config: config,
		log:    log,
		ctx:    context.Background(),
	}

}

func (h *Handler) Gather() {
	h.log.Debugf("Running on Org: %s", h.config.SourceOrg.Name)
	//var orgInfo OrganisationInformation
	//h.log.Debugf("Gathering Org Details")
	//org, err := h.orgDetails(h.ctx, h.config.SourceOrg.Name)
	//if err != nil {
	//	panic(err)
	//}
	//orgInfo.Organisation = org
	//h.log.Debugf("Gathering Repo Details")
	//repos, err := h.orgRepos(h.ctx)
	//if err != nil {
	//	panic(err)
	//}
	//orgInfo.Repositories = repos
	//var teams []Team
	//for _, v := range h.teamCache {
	//	teams = append(teams, v)
	//}
	//
	//orgInfo.SeenTeams = teams
	//h.printJson(orgInfo, orgInfo.Organisation.Name)
}

func (h *Handler) Plan() {
	if (h.config.SingleRepository == "" && !h.config.AllRepositories) || (h.config.SingleRepository != "" && h.config.AllRepositories) {
		panic("Must set either --repository or --all-repositories")
	}
	fmt.Println(h.config.SourceOrg.Name)
	fmt.Println(h.config.DestinationOrg.Name)
	gHubHandler := ghub.NewHandler(h.ctx, h.config.GithubToken, h.config, h.log)
	sourceOrg, err := gHubHandler.OrgDetails(h.ctx, h.config.SourceOrg.Name, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(sourceOrg.Repositories)
	destOrg, err := gHubHandler.OrgDetails(h.ctx, h.config.DestinationOrg.Name, false)
	if err != nil {
		panic(err)
	}
	teamMigration, teamRenames, missingMembers := h.teamsToMigrate(sourceOrg, destOrg)
	repoMigration := h.ValidateReposToMigrate(sourceOrg, destOrg)

	// TODO: Check repo names are unique and if not raise as issue for manual fixing
	// TODO: generate plan of:
	// Teams to create, reuse
	// Members that are missing
	// Repositories that will be migrated
	// Repositories that have name clashes and need manual fixes.
	type Plan struct {
		SourceOrgName      string
		DestinationOrgName string
		MissingMembers     []ghub.Member
		Teams              TeamMigration
		TeamsRenameMapping map[string]string
		Repositories       RepoMigration
	}
	plan := Plan{
		Teams:              teamMigration,
		TeamsRenameMapping: teamRenames,
		MissingMembers:     missingMembers,
		Repositories:       repoMigration,
		SourceOrgName:      sourceOrg.Name,
		DestinationOrgName: destOrg.Name,
	}
	h.printJson(plan, fmt.Sprintf("plan-%s-%s", sourceOrg.Name, destOrg.Name))

}

//
//func (h *Handler) planSingleRepo(repoName string) []Repository {
//	repo, err := h.repoDetails(h.ctx, repoName)
//	if err != nil {
//		h.log.Panicf("Failed to get repo details for %s, error: %s", repoName, err.Error())
//		panic(err)
//	}
//	return []Repository{repo}
//}
//
//func (h *Handler) planAllRepos() []Repository {
//	repos, err := h.orgRepos(h.ctx)
//	if err != nil {
//		h.log.Panicf("Failed to get all repositories for %s, error: %s", h.config.SourceOrg.Name, err.Error())
//
//		panic(err)
//	}
//	return repos
//}
