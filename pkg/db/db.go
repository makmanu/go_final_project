package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)


var db *sql.DB




func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if install {
		schema := `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(64) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT "");
CREATE INDEX scheduler_date ON scheduler (date);`
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}
	
	return nil
}