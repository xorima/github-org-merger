package cmd

import (
	"fmt"

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
	},
}

func init() {
	migrateCmd.AddCommand(planCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// planCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
