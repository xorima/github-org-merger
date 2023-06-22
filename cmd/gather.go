/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/xorima/github-org-merger/internal/config"
	"github.com/xorima/github-org-merger/internal/merger"

	"github.com/spf13/cobra"
)

// gatherCmd represents the gather command
var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Gathers all information needed for the migration and preps files for you to update",
	Long:  `:shrug:`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gather called")
		h := merger.NewHandler(config.AppConfig)
		h.Handle()
	},
}

func init() {
	rootCmd.AddCommand(gatherCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// Set the SourceOrg in the config
	gatherCmd.Flags().StringVarP(&config.AppConfig.SourceOrg.Name, "source-org", "s", "", "The source org to migrate from")
	err := gatherCmd.MarkFlagRequired("source-org")
	if err != nil {
		panic(err)
	}
}
