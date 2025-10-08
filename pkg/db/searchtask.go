package db

import (
	"database/sql"
	"strings"
	"time"
)


func parseDate(s string) (bool, string) {
	s = strings.TrimSpace(s)
	
	// Пытаемся распарсить дату
	t, err := time.Parse("02.01.2006", s)
	if err != nil {
		return false, ""
	}
	
	// Возвращаем дату в формате 20060102
	return true, t.Format("20060102")
}

func SearchTasks(search string, limit int) ([]*Task, error) {
	// Проверяем, является ли search датой в формате 02.01.2006
	if isDate, date := parseDate(search); isDate {
		// Поиск по дате с именованными параметрами
		rows, err := db.Query(`
			SELECT * FROM scheduler 
			WHERE date = :date 
			ORDER BY date 
			LIMIT :limit`,
			sql.Named("date", date),
			sql.Named("limit", limit))
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		return scanTasks(rows)
	}

	// Поиск по заголовку или комментарию с именованными параметрами
	searchPattern := "%" + search + "%"
	rows, err := db.Query(`
		SELECT * FROM scheduler 
		WHERE title LIKE :search OR comment LIKE :search 
		ORDER BY date 
		LIMIT :limit`,
		sql.Named("search", searchPattern),
		sql.Named("limit", limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return scanTasks(rows)
}