/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
  "nviro/cmd"
  "nviro/db"
)

func main() {
  db.OpenDatabase()
	cmd.Execute()
}
