package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
  "database/sql"
  "fmt"
  "io"
  "log"
  "os"
  "runtime"

  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDatabase() error {
  var err error
  fileName := "db.sqlite3"
  dbPath := "./storage/"
  homeDirectory, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
  osType := runtime.GOOS
  switch osType {
  case "windows":
    dbPath = homeDirectory + "/.local/share/nviro/"
  case "darwin":
    dbPath = homeDirectory + "/.local/share/nviro/"
  case "linux":
    dbPath = homeDirectory + "/.local/share/nviro/"
  default:
    fmt.Printf("%s.\n", osType)
  }
  _ , error := os.Stat(dbPath)

  if os.IsNotExist(error) {
    fmt.Printf("Creating dir: %s\n", dbPath)
  }
  {
    err := os.MkdirAll(dbPath, 0750)
    if err != nil && !os.IsExist(err) {
      log.Fatal(err)
    }
  }
  {
    filePath := dbPath + fileName
    db, err = sql.Open("sqlite3", filePath)
    if err != nil {
      return err
    }
  }
  return db.Ping()
}

func SchemaCreate() error {
  {
    createProjectTableSQL := `
      CREATE TABLE IF NOT EXISTS "projects"
      (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "name" VARCHAR NOT NULL,
        "created_at" DATETIME,
        "updated_at" DATETIME
      );`

    statement, err := db.Prepare(createProjectTableSQL)
    if err != nil {
      log.Fatal(err.Error())
    }
    statement.Exec()
    log.Println("Created projects table")
  }

  {
    createEnvironmentTableSQL := `
      CREATE TABLE IF NOT EXISTS "environments"
      (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "name" VARCHAR NOT NULL,
        "content" TEXT NOT NULL,
        "deleted_at" DATETIME,
        "project_id" INTEGER NOT NULL,
        "created_at" DATETIME,
        "updated_at" DATETIME,
        foreign key("project_id") REFERENCES "projects"("id")
      );`

    statement, err := db.Prepare(createEnvironmentTableSQL)
    if err != nil {
      log.Fatalln(err.Error())
    }
    statement.Exec()
    log.Println("Created environments table")
  }
  return nil
}

type Project struct {
	Id int `json:"id"`
	Name string `json:"name"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
}

func ProjectInsert(name string) {
  statement, err := db.Prepare(`
    INSERT INTO projects
      (name, created_at, updated_at)
    VALUES
      (?, datetime('now'), datetime('now'))
  `)
  if err != nil {
    log.Fatal(err.Error())
  }
  {
    _, err := statement.Exec(name)
    if err != nil {
      log.Fatalln(err.Error())
    }
  }
}

func ProjectList() []Project {
  rows, err := db.Query("SELECT * FROM projects ORDER BY name")
  if err != nil {
    log.Fatal(err.Error())
  }
  defer rows.Close()
  projects := []Project{}
  for rows.Next() {
    var id int
    var name string
    var created_at sql.NullString
    var updated_at sql.NullString
    err = rows.Scan(&id, &name, &created_at, &updated_at)
    if err != nil {
      log.Fatal(err.Error())
    }
    projects = append(projects, Project{
      id,
      name,
      created_at.String,
      updated_at.String,
    })
  }

  return projects
}

func ProjectDelete(id string) {
  statement, err := db.Prepare("DELETE FROM projects WHERE id = ?")
  if err != nil {
    log.Fatal(err.Error())
  }
  {
    _, err := statement.Exec(id)
    if err != nil {
      log.Fatalln(err.Error())
    }
  }
}

func EnvironmentInsert(name string, content string, project_id string) {
  contentEncrypted := encrypt(getKey(), content)
  statement, err := db.Prepare(`
    INSERT INTO environments
      (name, content, project_id, created_at, updated_at, deleted_at)
    VALUES
      (?, ?, ?, datetime('now'), datetime('now'), NULL)`)
  if err != nil {
    log.Fatal(err.Error())
  }
  {
    _, err := statement.Exec(name, contentEncrypted, project_id)
    if err != nil {
      log.Fatalln(err.Error())
    }
  }
}

type Environment struct {
	Id int `json:"id"`
	Name string `json:"name"`
  Content string `json:"content"`
  ProjectId int `json:"project_id"`
  ProjectName string `json:"project_name"`
  DeletedAt string `json:"deleted_at"`
  CreatedAt string `json:"created_at"`
  UpdatedAt string `json:"updated_at"`
}

func EnvironmentList(withDeleted bool) []Environment {
  whereSQL := "WHERE environments.deleted_at IS NULL"
  if withDeleted {
    whereSQL = ""
  }
  query := fmt.Sprintf(`
    SELECT
      environments.id,
      environments.name,
      environments.content,
      environments.project_id,
      projects.name AS project_name,
      environments.deleted_at,
      environments.created_at,
      environments.updated_at
    FROM
      environments
      INNER JOIN projects ON projects.id = environments.project_id
    %s
    ORDER BY environments.id ASC
  `, whereSQL)
  rows, err := db.Query(query)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer rows.Close()
  environments := []Environment{}
  key := getKey()
  for rows.Next() {
    var id int
    var name string
    var content string
    var project_id int
    var project_name string
    var deleted_at sql.NullString
    var created_at sql.NullString 
    var updated_at sql.NullString
    err = rows.Scan(&id, &name, &content, &project_id, &project_name, &deleted_at, &created_at, &updated_at)
    if err != nil {
      log.Fatal(err.Error())
    }
    contentDecrypted := decrypt(key, content)
    environments = append(environments, Environment{
      id,
      name,
      contentDecrypted,
      project_id,
      project_name,
      deleted_at.String,
      created_at.String,
      updated_at.String,
    })
  }

  return environments
}

func EnvironmentDelete(id string) {
  statement, err := db.Prepare("UPDATE environments SET deleted_at = datetime('now') WHERE id = ?")
  if err != nil {
    log.Fatal(err.Error())
  }
  {
    _, err := statement.Exec(id)
    if err != nil {
      log.Fatalln(err.Error())
    }
  }
}


func EnvironmentShow(id int) Environment {
  statement, err := db.Prepare(`
    SELECT
      environments.name,
      environments.content,
      environments.project_id,
      projects.name AS project_name,
      environments.deleted_at,
      environments.created_at,
      environments.updated_at
    FROM
      environments
      INNER JOIN projects ON projects.id = environments.project_id
    WHERE environments.id = ?
  `)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer statement.Close()
  var name string
  var content string
  var project_id int
  var project_name string
  var deleted_at sql.NullString
  var created_at sql.NullString 
  var updated_at sql.NullString
  {
    err = statement.QueryRow(id).Scan(&name, &content, &project_id, &project_name, &deleted_at, &created_at, &updated_at)
    if err != nil {
      log.Fatal(err.Error())
    }
  }
  key := getKey()
  contentDecrypted := decrypt(key, content)
  return Environment{
    id,
    name,
    contentDecrypted,
    project_id,
    project_name,
    deleted_at.String,
    created_at.String,
    updated_at.String,
  }
}

func encrypt(keyString string, stringToEncrypt string) (encryptedString string) {
	// convert key to bytes
	key, _    := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(keyString string, stringToDecrypt string) string {
	key, _ := hex.DecodeString(keyString)
	ciphertext, _ := base64.URLEncoding.DecodeString(stringToDecrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}

func getKey() string {
	key := []byte("T4rukYC8g5b9DkcbLuxxByCRM9hsrgN7")
	keyStr := hex.EncodeToString(key) //convert to string for saving
  return keyStr
}
