/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
  "envtool/cmd"
  "envtool/db"
)

func main() {
  db.OpenDatabase()
	cmd.Execute()
}
