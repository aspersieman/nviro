/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "fmt"
  "os"
  "text/tabwriter"
  "strconv"

	"github.com/spf13/cobra"
  "github.com/manifoldco/promptui"
  "nviro/db"
)

// environmentShowCmd represents the environmentShow command
var environmentShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show specific environment file",
	Long:  "Show specific environment file",
	Run: func(cmd *cobra.Command, args []string) {
		environmentShow()
	},
}

func init() {
	environmentCmd.AddCommand(environmentShowCmd)
}

func promptGetEnvironmentShowInput(pc promptContent) int {
  validate := func(input string) error {
    intInput, err := strconv.Atoi(input) 
    if err != nil {
      intInput = -1
    }
    if intInput < 0 {
      return fmt.Errorf(pc.errorMsg)
    }
    return nil
  }
  templates := &promptui.PromptTemplates{
    Prompt:  "{{ . }} ",
    Valid:   "{{ . | green }} ",
    Invalid: "{{ . | red }} ",
    Success: "{{ . | bold }} ",
  }
  prompt := promptui.Prompt{
    Label:     pc.label,
    Templates: templates,
    Validate:  validate,
  }

  result, err := prompt.Run()
  if err != nil {
    fmt.Printf("Prompt failed %v\n", err)
    os.Exit(1)
  }
  resultInt, err := strconv.Atoi(result) 
  if err != nil {
    resultInt = -1
  }

  return resultInt 
}

func environmentShow() {
  environmentIdPromptContent := promptContent{
    "Environment id is required",
    "Environment id",
  }

  environmentId := promptGetEnvironmentShowInput(environmentIdPromptContent)
  environment := db.EnvironmentShow(environmentId)
  
  w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
  fmt.Fprintln(w, "ID:\t", environment.Id)
  fmt.Fprintln(w, "Name:\t", environment.Name)
  fmt.Fprintln(w, "Project ID:\t", environment.ProjectId)
  fmt.Fprintln(w, "Project Name:\t", environment.ProjectName)
  fmt.Fprintln(w, "Deleted at:\t", environment.DeletedAt)
  fmt.Fprintln(w, "Created at:\t", environment.CreatedAt)
  fmt.Fprintln(w, "Updated at:\t", environment.UpdatedAt)
  fmt.Fprintln(w, "Content: \t")
  fmt.Fprintln(w, environment.Content)
  w.Flush()
}
