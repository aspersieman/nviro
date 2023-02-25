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

func promptGetEnvironmentShowInput(pc promptContent) string {
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

func environmentShow() {
  environmentIdPromptContent := promptContent{
    "Environment id is required",
    "Environment id",
  }

  environmentId := promptGetEnvironmentDeleteInput(environmentIdPromptContent)
  db.EnvironmentShow(environmentId)
}
