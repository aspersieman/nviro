/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  "github.com/manifoldco/promptui"
  "envtool/db"
)

var addCmd = &cobra.Command{
  Use:   "add",
  Short: "Add a new project",
  Long:  "Add a new project",
  Run: func(cmd *cobra.Command, args []string) {
    projectAdd()
  },
}

func init() {
  projectCmd.AddCommand(addCmd)
}

func promptGetProjectAddInput(pc promptContent) string {
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

func projectAdd() {
  projectNamePromptContent := promptContent{
    "Project name is required",
    "Project name",
  }

  projectName := promptGetProjectAddInput(projectNamePromptContent)
  db.ProjectInsert(projectName)
}
