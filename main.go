/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"nviro/cmd"
	"nviro/db"
)

func main() {
	fmt.Printf("STARTING APP")
	db.OpenDatabase()
	cmd.Execute()
}
