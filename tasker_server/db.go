package taskerserver

import (
	"regexp"
	"io/ioutil"
	"fmt"
	"database/sql"
	"path/filepath"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName = "tasker_server"
)


func init() {
	path, _ := filepath.Abs("./tasker_server/db/")

	if err := ReadAllSQLScripts(path, dbName); err != nil {
		panic(err)
	}
}
// ConnectToDB return poinetr to DB
func ConnectToDB() (*sql.DB, error) {
	return sql.Open("sqlite3", "./db/"+dbName+".db")
}

// ReadSQLScript MakeDBFile from sql script
func ReadSQLScript(filename , dbName string) error {
	sqlPath := fmt.Sprintf(filename)
	data, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		return err
	}

	pathToDBFile := fmt.Sprintf("./db/%s.db", dbName)
	db, err := sql.Open("sqlite3", pathToDBFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(string(data)); err != nil {
		return err
	}

	return nil
}
// ReadAllSQLScripts ....
func ReadAllSQLScripts(path ,dbName string) error {
	var scripts []string

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if match, _ := regexp.MatchString(`.*\.sql$`, file.Name()); match {
			scripts = append(scripts, path+"/"+file.Name())
		}
	}

	for _, script := range scripts {
		if err := ReadSQLScript(script, dbName); err != nil {
			return err
		}
	}

	return nil
}

//TODO рефактор