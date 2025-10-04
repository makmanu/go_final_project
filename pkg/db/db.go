package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)


var db *sql.DB

type Task struct {
    ID      string `json:"id"`
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

func Tasks(amount int64) ([]*Task, error){
	responseTasks := []*Task{}
	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY id LIMIT :amount", sql.Named("amount", amount))
	if err != nil {
		return responseTasks, err
	}
	defer rows.Close()

	for rows.Next() {
		var currentTask Task
	
		err := rows.Scan(&currentTask.ID, &currentTask.Date, &currentTask.Title, &currentTask.Comment, &currentTask.Repeat)
		if err != nil {
			return responseTasks, err
		}
		responseTasks = append(responseTasks, &currentTask)
	}
	return responseTasks, nil
}