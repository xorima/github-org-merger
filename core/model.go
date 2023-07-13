package core

type OrganizationQuery struct {
	Organization struct {
		Repositories struct {
			nodes []struct {
				Name  string
				Owner struct {
					Login string
				} `graphql:"owner"`

				//BranchProtectionRules struct {
				//	Nodes []struct {
				//		id                       string
				//		pattern                  string
				//		requiresApprovingReviews bool
				//		requiresCodeOwnerReviews bool
				//		requireLastPushApproval  bool
				//		dismissesStaleReviews    bool
				//		isAdminEnforced          bool
				//	} `graphql:"nodes"`
				//} `graphql:"branchProtectionRules(first:100)"`
			} `graphql:"nodes"`
		} `graphql:"repositories(first:$first, after:$cursor)"`
	} `graphql:"organization(login: $login)"`
}
