package cmd

import (
	"fmt"
	"github.com/xorima/github-org-merger/internal/config"
	"github.com/xorima/github-org-merger/internal/merger"

	"github.com/spf13/cobra"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Generate a plan for the migration to happen",
	Long: `A plan of everything to migrate will be generated. This will include:
	- Repositories
	- Teams
	- Users who are not in the new org
	- Branch Protection Rules to update
	- CodeOwners to update
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plan called")
		h := merger.NewHandler(config.AppConfig)
		h.Plan()
	},
}

func init() {
	migrateCmd.AddCommand(planCmd)
	migrateFlags(planCmd)
}
