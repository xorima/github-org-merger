package core

import (
	"context"
	"github.com/shurcooL/githubv4"
)

// gather will itterate through and find all repos with some form of branch protection and return them

func (h *Handler) Gather() {
	var query OrganizationQuery
	variables := map[string]interface{}{
		"first":  githubv4.Int(100),
		"login":  githubv4.String("summine"),
		"cursor": (*githubv4.String)(nil),
	}
	err := h.client.Query(context.Background(), &query, variables)
	if err != nil {
		panic(err)
	}
	//for _, r := range query.Organization.Repositories {
	//	fmt.Println(r.Name)
	//}

}
