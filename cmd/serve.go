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
  "mime"
  "net/http"
  "os"
  "path/filepath"
  "reflect"
  "regexp"
  "runtime"
  "strconv"
  "strings"
  "time"

	"github.com/spf13/cobra"
  "github.com/joho/godotenv"

  "nviro/db"
)

var (
  pages = map[string]string{
    "/": "static/templates/index.html",
  }
  //go:embed static/css/output.css static/dist/* static/templates/index.html static/img/* static/favicon.ico static/img/logo.png static/VERSION
  res embed.FS
  debug = true
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
	newRoute("GET",  "/", home),
	newRoute("GET",  "/static/.*", static),
	newRoute("GET",  "/static/favicon.ico", static),

	newRoute("GET",     "/api/projects", apiGetProjects),
	newRoute("POST",    "/api/projects", apiCreateProject),
	newRoute("PUT",     "/api/projects/([0-9]+)", apiUpdateProject),
	newRoute("DELETE",  "/api/projects/([0-9]+)", apiDeleteProject),

	newRoute("GET",     "/api/environments", apiGetEnvironments),
	newRoute("GET",     "/api/environments?[a-zA-Z0-9_&]*", apiGetEnvironments),
	newRoute("GET",     "/api/environments/([0-9]+)", apiShowEnvironment),
	newRoute("POST",    "/api/environments", apiCreateEnvironment),
	newRoute("PUT",     "/api/environments/([0-9]+)", apiUpdateEnvironment),
	newRoute("DELETE",  "/api/environments/([0-9]+)", apiDeleteEnvironment),
	newRoute("DELETE",  "/api/environments/([0-9]+)?[a-zA-Z0-9_&]*", apiDeleteEnvironment),
}

type ctxKey struct{}

type ProjectUpdate struct {
	Name string `json:"name"`
}

type EnvironmentUpdate struct {
	Name      string `json:"name"`
  Content   string `json:"content"`
  ProjectId int    `json:"project_id"`
}

type EnvironmentListParameters struct {
  Deleted bool    `json:"deleted"`
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
  versionFile, err := res.ReadFile("static/VERSION")
  if err != nil {
    log.Printf("ERROR: Failed to read version file: %s", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  version := string(versionFile)
  data := map[string]interface{}{
    "version": version,
  }
  if err := tpl.Execute(w, data); err != nil {
    return
  }
  err = tpl.ExecuteTemplate(w, "layout", version)
  if err != nil {
    log.Printf("ERROR: Cannot parse template 'layout'. %s", err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}

func static(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path
  path = path[1:]
  file, err := res.ReadFile(path)
  if err != nil {
    log.Printf("ERROR: Cannot find static file %s. %s", file, err.Error())
    http.Error(w, http.StatusText(404), 404)
    return
  }
  ext := filepath.Ext(path)
  mimeType := mime.TypeByExtension(ext)
  log.Printf("Serving static file %s (%s)", path, mimeType)
  if mimeType == "" {
    mimeType = http.DetectContentType(file)
  }
  log.Printf("mimeType: %s", mimeType)
  w.Header().Set("Content-Type", mimeType)
  _, err = w.Write(file)
  if err != nil {
    log.Printf("ERROR: Cannot write response. %s", err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}

func apiGetProjects(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  projects := db.ProjectList()
  err := json.NewEncoder(w).Encode(projects)
  if err != nil {
    log.Printf("ERROR: Cannot encode response. %s", err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
}

func apiCreateProject(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatal(err)
  }
  var project ProjectUpdate
  err = json.Unmarshal([]byte(reqBody), &project)
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
  if err != nil {
    log.Fatal(err)
  }
  var project ProjectUpdate
  err = json.Unmarshal([]byte(reqBody), &project)
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
  deleted := r.URL.Query().Get("deleted")
  del, _ := strconv.ParseBool(deleted)
  name := r.URL.Query().Get("name")
  project_id := r.URL.Query().Get("project_id")
  proj, _ := strconv.Atoi(project_id)
  idParam := r.URL.Query().Get("id")
  id, _ := strconv.Atoi(idParam)
  environments := db.EnvironmentList(del, name, proj, id)
  err := json.NewEncoder(w).Encode(environments)
  if err != nil {
    log.Fatal(err)
  }
}

func apiCreateEnvironment(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatal(err)
  }
  var environment EnvironmentUpdate
  err = json.Unmarshal([]byte(reqBody), &environment)
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
  id, _ := strconv.Atoi(getField(r, 0))
  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Fatal(err)
  }
  var environment EnvironmentUpdate
  err = json.Unmarshal([]byte(reqBody), &environment)
  if err != nil {
    log.Fatal(err)
  }
  err = db.EnvironmentUpdate(id, environment.Name, environment.Content, environment.ProjectId)
  if err != nil {
    log.Fatal(err)
  }
}

func apiDeleteEnvironment(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  id, _ := strconv.Atoi(getField(r, 0))
  forceQuery := r.URL.Query().Get("force")
  force, _ := strconv.ParseBool(forceQuery)
  err := db.EnvironmentDelete(id, force)
  if err != nil {
    log.Fatal(err)
  }
}

func apiShowEnvironment(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  id, _ := strconv.Atoi(getField(r, 0))
  environment := db.EnvironmentShow(id)
  err := json.NewEncoder(w).Encode(environment)
  if err != nil {
    log.Fatal(err)
  }
}

func loggingMiddleware(next http.HandlerFunc, method string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    f := getFunctionName(next)
    defer timeTrack(time.Now(), f)
    log.Printf("INFO: Request: [%s] %s => %s", method, r.URL.Path, f)
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
  p := "6969"
  log.Printf("INFO: Server started at :%s\n", p)
  http.Handle("favicon.ico", http.FileServer(http.FS(res)))
  http.Handle("static/", http.FileServer(http.FS(res)))
  http.Handle("static/img/", http.FileServer(http.FS(res)))
  http.Handle("static/js/", http.FileServer(http.FS(res)))
  http.Handle("static/css/", http.FileServer(http.FS(res)))
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

func getFunctionName(temp interface{}) string {
  strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
  return strs[len(strs)-1]
}

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  log.Printf("INFO: %s took %s", name, elapsed)
}

func goDotEnvVariable(key string) string {
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
