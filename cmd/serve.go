/*
Copyright © 2023 Nicol van der Merwe <aspersieman@gmail.com>

*/
package cmd

import (
  "embed"
  "encoding/json"
  "fmt"
  "html/template"
  "io/fs"
  "io/ioutil"
  "log"
  "net/http"

	"github.com/spf13/cobra"

  "nviro/db"
)

var (
  pages = map[string]string{
    "/": "static/templates/index.html",
  }
  //go:embed static/css/bootstrap/mixins/* static/css/bootstrap/utilities/* static/css/style.css static/js/* static/scss/bootstrap/mixins/* static/scss/bootstrap/utilities/* static/scss/bootstrap/vendor/* static/templates/index.html static/img/* static/favicon.ico 
  res embed.FS
  debug = false
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a web-based frontend for nviro",
	Long:  "Serve a web-based frontend for nviro",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
  if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
    if d.IsDir() {
      return nil
    }

    files = append(files, path)

    return nil
  }); err != nil {
    return nil, err
  }

  return files, nil
}

type ProjectAdd struct {
	Name string `json:"name"`
}

func projects(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  switch r.Method {
  case "GET":
    projects := db.ProjectList()
    json.NewEncoder(w).Encode(projects)
  case "POST":
    reqBody, err := ioutil.ReadAll(r.Body)
    var project ProjectAdd
    json.Unmarshal([]byte(reqBody), &project)
    if err != nil {
      log.Fatal(err)
    }
    err = db.ProjectInsert(project.Name)
    if err != nil {
      log.Fatal(err)
    }
  case "PUT":
    // update project
  default:
    w.WriteHeader(http.StatusNotImplemented)
    w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
  }  
}

type EnvironmentAdd struct {
	Name string `json:"name"`
  Content string `json:"content"`
  ProjectId int `json:"project_id"`
}

func environments(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")

  switch r.Method {
  case "GET":
    environments := db.EnvironmentList(true)
    json.NewEncoder(w).Encode(environments)
  case "POST":
    reqBody, err := ioutil.ReadAll(r.Body)
    var environment EnvironmentAdd
    json.Unmarshal([]byte(reqBody), &environment)
    if err != nil {
      log.Fatal(err)
    }
    err = db.EnvironmentInsert(environment.Name, environment.Content, environment.ProjectId)
    if err != nil {
      log.Fatal(err)
    }
  case "PUT":
    // update environment
  case "DELETE":
    // delete environment
  default:
    w.WriteHeader(http.StatusNotImplemented)
    w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
  }  
}

func serve() {
  files, _ := getAllFilenames(&res)
  if debug {
    fmt.Printf("INFO: Files in binary\n")
    for _, f := range files {
      fmt.Printf("INFO:\t- %s\n", f) 
    }
  }
  fmt.Println()
  http.Handle("/static/", http.FileServer(http.FS(res)))
  http.Handle("/favicon.ico", http.FileServer(http.FS(res)))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    page, ok := pages[r.URL.Path]
    fmt.Printf("Requested file: %s\n", page)
    if !ok {
      w.WriteHeader(http.StatusNotFound)
      return
    }
    tpl, err := template.ParseFS(res, page)
    if err != nil {
      log.Printf("ERROR: Page %s not found in pages cache...", r.RequestURI)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    w.Header().Set("Content-Type", "text/html")
    w.WriteHeader(http.StatusOK)
    data := map[string]interface{}{
      "userAgent": r.UserAgent(),
    }
    if err := tpl.Execute(w, data); err != nil {
      return
    }
    environments := db.EnvironmentList(true)
    err = tpl.ExecuteTemplate(w, "layout", environments)
    if err != nil {
      log.Print(err.Error())
      http.Error(w, http.StatusText(500), 500)
    }
  })
	http.HandleFunc("/api/projects", projects)
	http.HandleFunc("/api/environments", environments)
  p := "6969"
  log.Printf("Server started at :%s\n", p)
  err := http.ListenAndServe(":" + p, nil)
  if err != nil {
    panic(err)
  }
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
