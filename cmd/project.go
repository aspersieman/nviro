/*
Copyright © 2023 Nicol van der Merwe <aspersieman@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

type promptContent struct {
  errorMsg string
  label    string
}

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Add, delete, list or show projects.",
	Long: "Add, delete, list or show projects.",
}

func init() {
	rootCmd.AddCommand(projectCmd)
}
