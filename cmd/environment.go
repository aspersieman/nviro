/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// environmentCmd represents the environment command
var environmentCmd = &cobra.Command{
	Use:   "environment",
	Short: "Add, delete, list or show environments.",
	Long: "Add, delete, list or show environments.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("environment called")
	},
}

func init() {
	rootCmd.AddCommand(environmentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
