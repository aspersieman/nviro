/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nviro",
	Short: "Store and retrieve environment files",
	Long: `Save and store all your .env files in one place.
nviro will:
   - Store and encrypt your .env files
   - Allow you to retrieve your encrypted .env files
   - Keep track of changes to your env files`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
