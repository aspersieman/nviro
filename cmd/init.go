/*
Copyright © 2023 Nicol van der Merwe <aspersieman@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
  "nviro/db"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new nviro database and schema",
	Long: "Initialise a new nviro database and schema",
	Run: func(cmd *cobra.Command, args []string) {
		db.SchemaCreate()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
