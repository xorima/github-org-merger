/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xorima/github-org-merger/internal/config"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate the github repos from one org to another",
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	migrateFlags(migrateCmd)
}

func migrateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&config.AppConfig.DestinationOrg.Name, "destination-org", "d", "", "The destination org to migrate to")
	err := cmd.MarkFlagRequired("destination-org")
	if err != nil {
		panic(err)
	}
	//cmd.Flags().StringVarP(&config.AppConfig.SingleRepository, "repository", "r", "", "The single repository to migrate (optional, this or --all-repositories must be set)")
	cmd.Flags().BoolVarP(&config.AppConfig.AllRepositories, "all-repositories", "a", false, "Migrate all repositories (optional, this or --repository must be set)")

}
