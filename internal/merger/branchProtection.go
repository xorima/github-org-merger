package merger

import (
	"context"
	"fmt"
	"github.com/shurcooL/githubv4"
)

type Actor struct {
	Typename githubv4.String `graphql:"typename :__typename"`
	Team     struct {
		CombinedSlug githubv4.String `graphql:"combinedSlug"`
	} `graphql:"... on Team"`
}

type BypassPullRequestAllowances struct {
	Nodes []struct {
		Actor Actor
	} `graphql:"nodes"`
}

type BypassForcePushAllowances struct {
	Nodes []struct {
		Actor Actor
	} `graphql:"nodes"`
}

type BranchProtectionRule struct {
	Pattern                     githubv4.String
	BypassPullRequestAllowances BypassPullRequestAllowances `graphql:"bypassPullRequestAllowances(first: 100)"`
	BypassForcePushAllowances   BypassForcePushAllowances   `graphql:"bypassForcePushAllowances(first: 100)"`
}

type BranchProtectionRules struct {
	Nodes []BranchProtectionRule `graphql:"nodes"`
}

type GraphQLRepository struct {
	BranchProtectionRules BranchProtectionRules `graphql:"branchProtectionRules(first: 100)"`
}

type ProtectionBypasses struct {
	BranchPattern          string
	BypassPullRequestTeams []string
	BypassForcePush        []string
}

func (h *Handler) getBranchProtectionGroups(ctx context.Context, owner, repo string) ([]ProtectionBypasses, error) {

	type Query struct {
		Repository GraphQLRepository `graphql:"repository(owner: $owner, name: $repository)"`
	}
	var query Query
	variables := map[string]interface{}{
		"repository": githubv4.String(repo),
		"owner":      githubv4.String(owner),
	}
	err := h.clientGraphql.Query(ctx, &query, variables)
	if err != nil {
		h.log.Errorf("Unable to get branch protection rules for repo %s, error: %s", repo, err.Error())
		return nil, err
	}

	var results []ProtectionBypasses
	fmt.Println(len(query.Repository.BranchProtectionRules.Nodes))
	for _, node := range query.Repository.BranchProtectionRules.Nodes {
		if (len(node.BypassPullRequestAllowances.Nodes) > 0) || (len(node.BypassForcePushAllowances.Nodes) > 0) {
			tmp := ProtectionBypasses{
				BranchPattern: string(node.Pattern),
			}
			for _, t := range node.BypassPullRequestAllowances.Nodes {
				tmp.BypassPullRequestTeams = append(tmp.BypassPullRequestTeams, string(t.Actor.Team.CombinedSlug))
			}
			for _, t := range node.BypassForcePushAllowances.Nodes {
				tmp.BypassPullRequestTeams = append(tmp.BypassForcePush, string(t.Actor.Team.CombinedSlug))
			}
			results = append(results, tmp)
		}
	}
	return results, nil
}
