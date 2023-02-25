/*
Copyright © 2023 Nicol van der Merwe <aspersieman@gmail.com>

*/
package cmd

import (
	"github.com/spf13/cobra"
  "envtool/db"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialise a new envtool database and schema",
	Long: "Initialise a new envtool database and schema",
	Run: func(cmd *cobra.Command, args []string) {
		db.SchemaCreate()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
