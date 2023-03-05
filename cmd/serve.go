/*
Copyright © 2023 Nicol van der Merwe <aspersieman@gmail.com>

*/
package cmd

import (
  "context"
  "embed"
  "encoding/json"
  "fmt"
  "html/template"
  "io/fs"
  "io/ioutil"
  "log"
  "net/http"
  "reflect"
  "regexp"
  "runtime"
  "strconv"
  "strings"

	"github.com/spf13/cobra"

  "nviro/db"
)

var (
  pages = map[string]string{
    "/": "static/templates/index.html",
  }
  //go:embed static/css/style.css static/js/* static/templates/index.html static/img/* static/favicon.ico 
  res embed.FS
  debug = false
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a web-based frontend for nviro",
	Long:  "Serve a web-based frontend for nviro",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

var routes = []route{
	newRoute("GET", "/", home),

	newRoute("GET", "/api/projects", apiGetProjects),
	newRoute("POST", "/api/projects", apiCreateProject),
	newRoute("PUT", "/api/projects/([0-9]+)", apiUpdateProject),
	newRoute("DELETE", "/api/projects/([0-9]+)", apiDeleteProject),

	newRoute("GET", "/api/environments", apiGetEnvironments),
	newRoute("POST", "/api/environments", apiCreateEnvironment),
	newRoute("PUT", "/api/environments/([0-9]+)", apiUpdateEnvironment),
	newRoute("DELETE", "/api/environments/([0-9]+)", apiDeleteEnvironment),
}

type ctxKey struct{}

type ProjectUpdate struct {
	Name string `json:"name"`
}

type EnvironmentUpdate struct {
	Name string `json:"name"`
  Content string `json:"content"`
  ProjectId int `json:"project_id"`
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{
    method,
    regexp.MustCompile("^" + pattern + "$"),
    loggingMiddleware(http.HandlerFunc(handler), method),
  }
}

func home(w http.ResponseWriter, r *http.Request) {
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
}

func apiGetProjects(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  environments := db.ProjectList()
  json.NewEncoder(w).Encode(environments)
}

func apiCreateProject(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  reqBody, err := ioutil.ReadAll(r.Body)
  var project ProjectUpdate
  json.Unmarshal([]byte(reqBody), &project)
  if err != nil {
    log.Fatal(err)
  }
  err = db.ProjectInsert(project.Name)
  if err != nil {
    log.Fatal(err)
  }
}

func apiUpdateProject(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  reqBody, err := ioutil.ReadAll(r.Body)
  var project ProjectUpdate
  json.Unmarshal([]byte(reqBody), &project)
  if err != nil {
    log.Fatal(err)
  }
  id, _ := strconv.Atoi(getField(r, 0))
  err = db.ProjectUpdate(id, project.Name)
  if err != nil {
    log.Fatal(err)
  }
}

func apiDeleteProject(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  id, _ := strconv.Atoi(getField(r, 0))
  err := db.ProjectDelete(id)
  if err != nil {
    log.Fatal(err)
  }
}

func apiGetEnvironments(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  environments := db.EnvironmentList(true)
  json.NewEncoder(w).Encode(environments)
}

func apiCreateEnvironment(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  reqBody, err := ioutil.ReadAll(r.Body)
  var environment EnvironmentUpdate
  json.Unmarshal([]byte(reqBody), &environment)
  if err != nil {
    log.Fatal(err)
  }
  err = db.EnvironmentInsert(environment.Name, environment.Content, environment.ProjectId)
  if err != nil {
    log.Fatal(err)
  }
}

func apiUpdateEnvironment(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  reqBody, err := ioutil.ReadAll(r.Body)
  var environment EnvironmentUpdate
  json.Unmarshal([]byte(reqBody), &environment)
  if err != nil {
    log.Fatal(err)
  }
  id, _ := strconv.Atoi(getField(r, 0))
  err = db.EnvironmentUpdate(id, environment.Name, environment.Content, environment.ProjectId)
  if err != nil {
    log.Fatal(err)
  }
}

func apiDeleteEnvironment(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  id, _ := strconv.Atoi(getField(r, 0))
  err := db.EnvironmentDelete(id)
  if err != nil {
    log.Fatal(err)
  }
}

func getFunctionName(temp interface{}) string {
  strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
  return strs[len(strs)-1]
}

func loggingMiddleware(next http.HandlerFunc, method string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("INFO: Request: [%s] %s => %s", method, r.URL.Path, getFunctionName(next))
		next.ServeHTTP(w, r)
	})
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}
			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func start() {
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
  p := "6969"
  log.Printf("INFO: Server started at :%s\n", p)
  err := http.ListenAndServe(":" + p, http.HandlerFunc(Serve))
  if err != nil {
    panic(err)
  }
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

func init() {
	rootCmd.AddCommand(serveCmd)
}
