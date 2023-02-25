package db

import (
  "fmt"
  "database/sql"
  "log"
  "os"
  "io"
  "text/tabwriter"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"

  _ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func OpenDatabase() error {
  var err error
  // TODO: use goos to detect the OS: https://pkg.go.dev/runtime#GOOS
  //    Use this location to place the sqlite3.db file in the appropriate folder
  path := "./storage/db.sqlite3"
  db, err = sql.Open("sqlite3", path)
  if err != nil {
    return err
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

func ProjectList() {
  // TODO handle presentation logic outside of this
  rows, err := db.Query("SELECT * FROM projects ORDER BY name")
  if err != nil {
    log.Fatal(err.Error())
  }
  defer rows.Close()
  w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
  fmt.Fprintln(w, "ID\t", "NAME")
  fmt.Fprintln(w, "----\t", "-----")
  for rows.Next() {
    var id int
    var name string
    var created_at sql.NullString
    var updated_at sql.NullString
    err = rows.Scan(&id, &name, &created_at, &updated_at)
    if err != nil {
      log.Fatal(err.Error())
    }
    fmt.Fprintln(w, fmt.Sprintf("%d\t%s", id, name))
    w.Flush()
  }
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

func EnvironmentList(withDeleted bool) {
  // TODO handle presentation logic outside of this
  whereSQL := ""
  if withDeleted {
    whereSQL = "WHERE environments.deleted_at IS NOT NULL"
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
  w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
  // TODO use a slice to define column headers and dividers
  columnHeaders := "ID\tNAME\tPROJECT ID\tPROJECT NAME\tDELETED AT\tCREATED AT\tUPDATED AT"
  columnDivider := "--\t----\t----------\t------------\t\t----------\t----------"
  if !withDeleted {
    columnHeaders = "ID\tNAME\tPROJECT ID\tPROJECT NAME\tCREATED AT\tUPDATED AT"
    columnDivider = "--\t----\t----------\t------------\t----------\t----------"
  }
  fmt.Fprintln(w, columnHeaders)
  fmt.Fprintln(w, columnDivider)
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
    if (withDeleted) {
      data := fmt.Sprintf("%d\t%s\t%d\t%s\t%s\t%s\t%s", id, name, project_id, project_name, deleted_at.String, created_at.String, updated_at.String)
      fmt.Fprintln(w, data)
    } else {
      data := fmt.Sprintf("%d\t%s\t%d\t%s\t%s\t%s", id, name, project_id, project_name, created_at.String, updated_at.String)
      fmt.Fprintln(w, data)
    }
  }
  w.Flush()
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


func EnvironmentShow(id string) {
  // TODO handle presentation logic outside of this
  statement, err := db.Prepare(`
    SELECT
      environments.content
    FROM
      environments
    WHERE id = ?
  `)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer statement.Close()
  var content string
  {
    err = statement.QueryRow(id).Scan(&content)
    if err != nil {
      log.Fatal(err.Error())
    }
  }
  fmt.Println("----------")
  key := getKey()
  contentDecrypted := decrypt(key, content)
  fmt.Println(contentDecrypted)
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
