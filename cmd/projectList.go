/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
  "envtool/db"
)

// projectListCmd represents the projectList command
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects.",
	Long:  "List all projects.",
	Run: func(cmd *cobra.Command, args []string) {
		db.ProjectList()
	},
}

func init() {
	projectCmd.AddCommand(projectListCmd)
}
