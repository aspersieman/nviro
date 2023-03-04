/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
  "fmt"
  "os"
  "strconv"

  "github.com/spf13/cobra"
  "github.com/manifoldco/promptui"
  "nviro/db"
)

// environmentDeleteCmd represents the environmentDelete command
var environmentDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment",
	Long:  "Delete an environment",
	Run: func(cmd *cobra.Command, args []string) {
		environmentDelete()
	},
}

func init() {
	environmentCmd.AddCommand(environmentDeleteCmd)
}

func promptGetEnvironmentDeleteInput(pc promptContent) string {
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

func environmentDelete() {
  environmentIdPromptContent := promptContent{
    "Environment id is required",
    "Environment id",
  }
  environmentId, _ := strconv.Atoi(promptGetEnvironmentDeleteInput(environmentIdPromptContent))
  db.EnvironmentDelete(environmentId)
  fmt.Printf("Environment %s deleted\n", environmentId)
}
