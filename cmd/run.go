/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/xorima/github-org-merger/internal/config"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the migration for the inputted plan",
	Long: `Executes the given plan and migrates resources within it. This will include:
	- Repositories
	- Teams
	- Users who are not in the new org
	- Branch Protection Rules to update
	- CodeOwners to update
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
}

func init() {
	migrateCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	runCmd.Flags().StringVarP(&config.AppConfig.PlanFile, "plan-file", "f", "", "The plan file to use for the migration")
	err := runCmd.MarkFlagRequired("plan-file")
	if err != nil {
		panic(err)
	}
	migrateFlags(runCmd)

}
