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

// environmentListCmd represents the environmentList command
var environmentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environments.",
	Long:  "List all environments.",
	Run: func(cmd *cobra.Command, args []string) {
    environmentList(true)
	},
}

func init() {
	environmentCmd.AddCommand(environmentListCmd)
}

func environmentList(withDeleted bool) {
  environments := db.EnvironmentList(withDeleted)

  w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
  // TODO use a slice to define column headers and dividers
  columnHeaders := "ID\tNAME\tPROJECT ID\tPROJECT NAME\tDELETED AT\tCREATED AT\tUPDATED AT"
  columnDivider := "--\t----\t----------\t------------\t----------\t----------\t----------"
  if !withDeleted {
    columnHeaders = "ID\tNAME\tPROJECT ID\tPROJECT NAME\tCREATED AT\tUPDATED AT"
    columnDivider = "--\t----\t----------\t------------\t----------\t----------"
  }
  fmt.Fprintln(w, columnHeaders)
  fmt.Fprintln(w, columnDivider)

	for _, environment := range environments {
    if (withDeleted) {
      data := fmt.Sprintf(
        "%d\t%s\t%d\t%s\t%s\t%s\t%s",
        environment.Id,
        environment.Name,
        environment.ProjectId,
        environment.ProjectName,
        environment.DeletedAt,
        environment.CreatedAt,
        environment.UpdatedAt,
      )
      fmt.Fprintln(w, data)
    } else {
      data := fmt.Sprintf(
        "%d\t%s\t%d\t%s\t%s\t%s",
        environment.Id,
        environment.Name,
        environment.ProjectId,
        environment.ProjectName,
        environment.CreatedAt,
        environment.UpdatedAt,
      )
      fmt.Fprintln(w, data)
    }
	}

  w.Flush()
}
