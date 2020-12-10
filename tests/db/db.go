package db

import (
	"regexp"
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// PathToDBDir path to this dir
const (
	PathToDBDir = "./db/"
)

func init() {
	data, err := ioutil.ReadFile("./db/go-to-do.sql")
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	db, err := sql.Open("sqlite3", "./db/go-to-do.db")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer db.Close()

	_, err = db.Exec(string(data))
	if err != nil {
		log.Fatal(err)
	}
}
// ConnectToDB return poinetr to DB
func ConnectToDB(db Databaser) (*sql.DB, error) {
	return sql.Open("sqlite3", "./db/"+db.DB()+".db")
}

// ReadSQLScript MakeDBFile from sql script
func ReadSQLScript(filename string, d Databaser) error {
	sqlPath := fmt.Sprintf(filename)
	data, err := ioutil.ReadFile(sqlPath)
	if err != nil {
		return err
	}

	pathToDBFile := fmt.Sprintf("./db/%s.db", d.DB())
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
func ReadAllSQLScripts(path string, d Databaser) error {
	var scripts []string

	files, _ := ioutil.ReadDir(path)
	for _, file := range files {
		if match, _ := regexp.MatchString(`.*\.sql$`, file.Name()); match {
			scripts = append(scripts, path+"/"+file.Name())
		}
	}

	for _, script := range scripts {
		if err := ReadSQLScript(script, d); err != nil {
			return err
		}
	}

	return nil
}