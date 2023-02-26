/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "text/tabwriter"
  "os"
  "fmt"

	"github.com/spf13/cobra"
  "nviro/db"
)

// projectListCmd represents the projectList command
var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects.",
	Long:  "List all projects.",
	Run: func(cmd *cobra.Command, args []string) {
		projectList()
	},
}

func init() {
	projectCmd.AddCommand(projectListCmd)
}

func projectList() {
  projects := db.ProjectList()

  w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
  fmt.Fprintln(w, "ID\tNAME\tCREATED AT\tUPDATED AT")
  fmt.Fprintln(w, "--\t----\t----------\t----------")
	for _, project := range projects {
    data := fmt.Sprintf(
      "%d\t%s\t%s\t%s",
      project.Id,
      project.Name,
      project.CreatedAt,
      project.UpdatedAt,
    )
    fmt.Fprintln(w, data)
  }
    w.Flush()
}
