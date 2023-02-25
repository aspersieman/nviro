/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
  "envtool/db"
)

// environmentListCmd represents the environmentList command
var environmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments.",
	Long:  "List all environments.",
	Run: func(cmd *cobra.Command, args []string) {
		db.EnvironmentList(false)
	},
}

func init() {
	environmentCmd.AddCommand(environmentListCmd)
}
