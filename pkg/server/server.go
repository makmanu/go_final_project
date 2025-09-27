package server

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Start_server(port string) error {
	webDir := "./web"
	log.Println("Запуск сервера")
	defer log.Println("Остановка сервера")

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(port, nil)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}
	
func afterNow(date time.Time, now time.Time) bool {
	dur := date.Sub(now)
	if dur > 0 {
		return true
	}
	return false
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	
	if len(repeat) == 0 {
		return "", errors.New("repeat rule is missing")
	}

	dateTime, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}
	repeatRules := strings.Split(repeat, " ")

	resultDate := dateTime
	switch repeatRules[0] {
	case "y":
		for !afterNow(resultDate, now) {
			resultDate = resultDate.AddDate(1, 0, 0)
		}
	case "d":
		numberOfDays, err := strconv.Atoi(repeatRules[1])
		if err != nil {
			return "", err
		}
		if numberOfDays < 1 || numberOfDays > 400 {
			return "", errors.New("day amount not in range [1,400]")
		}
		for !afterNow(resultDate, now) {
			resultDate = resultDate.AddDate(0, 0, numberOfDays)
		}
	default:
		return "", errors.New("invalid rule")
	}	

	return resultDate.Format("20060102"), nil
}