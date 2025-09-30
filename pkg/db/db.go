package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)


var db *sql.DB

type Task struct {
    ID      int64 `json:"id"`
    Date    string `json:"date"`
    Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}


func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

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

func AddTask(task *Task) (int64, error) {
    var id int64
    query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
    res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
    if err == nil {
        id, err = res.LastInsertId()
    }
    return id, err
} 