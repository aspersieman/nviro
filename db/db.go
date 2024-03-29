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
  "strings"

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
  EnvironmentCount int `json:"environment_count"`
}

func ProjectInsert(name string) error {
  statement, err := db.Prepare(`
    INSERT INTO projects
      (name, created_at, updated_at)
    VALUES
      (?, datetime('now'), datetime('now'))
  `)
  if err != nil {
    log.Fatal(err.Error())
    return err
  }
  _, err = statement.Exec(name)
  if err != nil {
    log.Fatalln(err.Error())
    return err
  }
  return nil
}

func ProjectUpdate(id int, name string) error {
  statement, err := db.Prepare(`
    UPDATE projects
    SET
      name = ?,
      updated_at = datetime('now')
    WHERE id = ?
  `)
  if err != nil {
    log.Fatal(err.Error())
    return err
  }
  _, err = statement.Exec(name, id)
  if err != nil {
    log.Fatalln(err.Error())
    return err
  }
  return nil
}

func ProjectList() []Project {
  rows, err := db.Query(`
    SELECT
      projects.*,
      (SELECT COUNT(*) FROM environments WHERE project_id = projects.id) AS environments_count
    FROM projects ORDER BY name ASC
  `)
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
    var environment_count int
    err = rows.Scan(&id, &name, &created_at, &updated_at, &environment_count)
    if err != nil {
      log.Fatal(err.Error())
    }
    projects = append(projects, Project{
      id,
      name,
      created_at.String,
      updated_at.String,
      environment_count,
    })
  }

  return projects
}

func ProjectDelete(id int) error {
  statement, err := db.Prepare("DELETE FROM projects WHERE id = ?")
  if err != nil {
    log.Fatal(err.Error())
  }
  _, err = statement.Exec(id)
  if err != nil {
    log.Fatalln(err.Error())
    return err
  }
  return nil
}

func EnvironmentInsert(name string, content string, project_id int) error {
  contentEncrypted := encrypt(getKey(), content)
  statement, err := db.Prepare(`
    INSERT INTO environments
      (name, content, project_id, created_at, updated_at, deleted_at)
    VALUES
      (?, ?, ?, datetime('now'), datetime('now'), NULL)`)
  if err != nil {
    log.Fatal(err.Error())
    return err
  }
  _, err = statement.Exec(name, contentEncrypted, project_id)
  if err != nil {
    log.Fatalln(err.Error())
    return err
  }
  return nil
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
  HistoryCount int `json:"history_count"`
}

func EnvironmentList(withDeleted bool, name string, project_id int, id int) []Environment {
  where := make([]string, 0)
  var parameters []interface{}
  whereSQL := ""
  if !withDeleted {
    where = append(where, "environments.deleted_at IS NULL")
  }
  if name != "" {
    where = append(where, "environments.name = ?")
    parameters = append(parameters, name)
  }
  if project_id > 0 {
    where = append(where, "environments.project_id = ?")
    parameters = append(parameters, project_id)
  }
  if id > 0 {
    where = append(where, "environments.id != ?")
    parameters = append(parameters, id)
  }
  if len(where) > 0 {
    whereSQL = "WHERE " + strings.Join(where, " AND ")  
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
      environments.updated_at,
      (
        SELECT
          COUNT(*)
        FROM
          environments history
        WHERE
          history.name = environments.name
          AND history.project_id = environments.project_id
          AND history.id != environments.id
      ) AS history_count
    FROM
      environments
      INNER JOIN projects ON projects.id = environments.project_id
    %s
    ORDER BY
      environments.name ASC,
      environments.project_id ASC,
      environments.created_at ASC
  `, whereSQL)
  rows, err := db.Query(query, parameters...)
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
    var history_count int
    err = rows.Scan(&id, &name, &content, &project_id, &project_name, &deleted_at, &created_at, &updated_at, &history_count)
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
      history_count,
    })
  }

  return environments
}

func EnvironmentDelete(id int, force bool) error {
  if !force {
    statement, err := db.Prepare(`
      SELECT
        environments.id,
        environments.deleted_at
      FROM
        environments
      WHERE
        environments.id = ?
    `)
    if err != nil {
      log.Fatal(err.Error())
    }
    defer statement.Close()
    var environmentId int
    var deleted_at sql.NullString
    err = statement.QueryRow(id).Scan(&environmentId, &deleted_at)
    if err != nil && err != sql.ErrNoRows {
      log.Fatal(err.Error())
    }
    if environmentId > 0 && !deleted_at.Valid {
      statement, err := db.Prepare("UPDATE environments SET deleted_at = datetime('now') WHERE id = ?")
      if err != nil {
        log.Fatal(err.Error())
        return err
      }
      _, err = statement.Exec(id)
      if err != nil {
        log.Fatalln(err.Error())
        return err
      }
    } 
  } else {
    statement, err := db.Prepare("DELETE FROM environments WHERE id = ?")
    if err != nil {
      log.Fatal(err.Error())
      return err
    }
    _, err = statement.Exec(id)
    if err != nil {
      log.Fatalln(err.Error())
      return err
    }
  }
  return nil
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
      environments.updated_at,
      (
        SELECT
          COUNT(*)
        FROM
          environments history
        WHERE
          history.name = environments.name
          AND history.project_id = environments.project_id
          AND history.id != environments.id
      ) AS history_count
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
  var history_count int
  err = statement.QueryRow(id).Scan(&name, &content, &project_id, &project_name, &deleted_at, &created_at, &updated_at, &history_count)
  if err != nil {
    log.Fatal(err.Error())
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
    history_count,
  }
}

func EnvironmentUpdate(id int, name string, content string, project_id int) error {
  statement, err := db.Prepare(`
    SELECT
      environments.id
    FROM
      environments
    WHERE
      environments.name = ?
      AND environments.project_id = ?
      AND environments.deleted_at IS NULL
  `)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer statement.Close()
  var environmentId int
  err = statement.QueryRow(name, project_id).Scan(&environmentId)
  if environmentId > 0 {
    errDelete := EnvironmentDelete(id, false)
    if errDelete != nil {
      log.Fatal(errDelete)
    }
    err := EnvironmentInsert(name, content, project_id)
    if err != nil {
      log.Fatal(err.Error())
      return err
    }
  } else {
    statement, err := db.Prepare(`
      UPDATE environments
      SET
        name = ?,
        content = ?,
        project_id = ?,
        updated_at = datetime('now')
      WHERE id = ?
    `)
    if err != nil {
      log.Fatal(err.Error())
      return err
    }
    contentEncrypted := encrypt(getKey(), content)
    _, err = statement.Exec(name, contentEncrypted, project_id, id)
    if err != nil {
      log.Fatalln(err.Error())
      return err
    }
  }
  return nil
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
