/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
  "log"
	"net/http"
  "encoding/json"
  "html/template"
  "path/filepath"
  "os"

	"github.com/spf13/cobra"

  "nviro/db"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a web-based frontend for nviro",
	Long: "Serve a web-based frontend for nviro",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func projects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

  projects := db.ProjectList()

  json.NewEncoder(w).Encode(projects)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("./static/templates", "index.html")
  fpEnd := filepath.Clean(r.URL.Path)
  if fpEnd == "/" {
    // Default to index.html
    fpEnd = "/index.html"
  }
	fp := filepath.Join("./static/templates", fpEnd)
  fmt.Printf("Requested file: %s\n", fpEnd)

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		// Log the detailed error
		log.Print(err.Error())
		// Return a generic "Internal Server Error" message
		http.Error(w, http.StatusText(500), 500)
		return
	}

  environments := db.EnvironmentList(true)

	err = tmpl.ExecuteTemplate(w, "layout", environments)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(500), 500)
	}
}

func serve() {
  fs := http.FileServer(http.Dir("./static"))
  http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", serveTemplate)
	http.HandleFunc("/api/projects", projects)

  p := "6969"
  fmt.Printf("Serving on port %s\n", p)
	log.Fatal(http.ListenAndServe(":" + p, nil))
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
