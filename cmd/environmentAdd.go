/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "fmt"
  "os"
  "io/ioutil"

  "github.com/spf13/cobra"
  "github.com/manifoldco/promptui"
  "envtool/db"
)

// environmentAddCmd represents the environmentAdd command
var environmentAddCmd = &cobra.Command{
  Use:   "add",
  Short: "Add a new environment",
  Long:  "Add a new environment",
  Run: func(cmd *cobra.Command, args []string) {
    environmentAdd()
  },
}

func init() {
  environmentCmd.AddCommand(environmentAddCmd)
}

// TODO: centralise this as it's used in multiple places
func promptGetEnvironmentAddInput(pc promptContent) string {
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

func environmentAdd() {
  environmentNamePromptContent := promptContent{
    "Environment name is required",
    "Environment name",
  }
  environmentContentPromptContent := promptContent{
    "Environment file is required",
    "Environment file",
  }
  environmentProjectIdPromptContent := promptContent{
    "Environment project ID is required",
    "Environment project ID",
  }

  environmentName := promptGetEnvironmentAddInput(environmentNamePromptContent)
  environmentContent := promptGetEnvironmentAddInput(environmentContentPromptContent)
  environmentProjectId := promptGetEnvironmentAddInput(environmentProjectIdPromptContent)

  if _, err := os.Stat(environmentContent); err != nil {
    fmt.Printf("File '%s' does not exist\n", environmentContent);
  } else {
    fileBytes, err := ioutil.ReadFile(environmentContent)
    if err != nil {
      panic(err)
    }
    fileString := string(fileBytes)

    db.EnvironmentInsert(environmentName, fileString, environmentProjectId)
  }
}
