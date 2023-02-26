/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/manifoldco/promptui"
  "nviro/db"
)

// projectDeleteCmd represents the projectDelete command
var projectDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a project",
	Long:  "Delete a project",
	Run: func(cmd *cobra.Command, args []string) {
		projectDelete()
	},
}

func init() {
	projectCmd.AddCommand(projectDeleteCmd)
}

func promptGetProjectDeleteInput(pc promptContent) string {
  validate := func(input string) error {
    if input == "" {
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

  return result
}

func projectDelete() {
  projectIdPromptContent := promptContent{
    "Project id is required",
    "Project id",
  }

  projectId := promptGetProjectDeleteInput(projectIdPromptContent)
  db.ProjectDelete(projectId)
  fmt.Printf("Project %s deleted\n", projectId)
}
