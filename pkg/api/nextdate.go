package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)


func HandleQuery(w http.ResponseWriter, req *http.Request) {
	var now string
	log.Println("start to handle /api/nextdate")
	if req.URL.Query().Has("now"){
		now = req.URL.Query().Get("now")
	} else {
		now = time.Now().Format("20060102")
	}
	if !(req.URL.Query().Has("date") && req.URL.Query().Has("repeat")) {
		fmt.Fprint(w, "err: not enough args\nwant to see[date, repeat]")
		return
	}
	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	response, err := nextDate(now, date, repeat)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprint(w, response)

} 

func afterNow(date time.Time, now time.Time) bool {
	dur := date.Sub(now)
	return dur > 0 
}

func nextDate(now string, dstart string, repeat string) (string, error) {

	if len(repeat) == 0 {
		return "", errors.New("repeat rule is missing")
	}

	nowParsed, err := time.Parse("20060102", now)
	if err != nil {
		return "", err
	}

	dateTime, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", err
	}
	repeatRules := strings.Split(repeat, " ")

	resultDate := dateTime
	switch repeatRules[0] {
	case "y":
		resultDate = resultDate.AddDate(1, 0, 0)
		for !afterNow(resultDate, nowParsed) {
			resultDate = resultDate.AddDate(1, 0, 0)
		}
	case "d":
		if len(repeatRules) != 2 {
			return "", errors.New("day repeat rule is invalid")
		}
		numberOfDays, err := strconv.Atoi(repeatRules[1])
		if err != nil {
			return "", err
		}
		if numberOfDays < 1 || numberOfDays > 400 {
			return "", errors.New("day amount not in range [1,400]")
		}
		resultDate = resultDate.AddDate(0, 0, numberOfDays)
		for !afterNow(resultDate, nowParsed) {
			resultDate = resultDate.AddDate(0, 0, numberOfDays)
		}
	default:
		return "", errors.New("invalid rule for next date: "+repeatRules[0])
	}	

	return resultDate.Format("20060102"), nil
}