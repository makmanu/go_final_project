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
		now = time.Now().Format(dateLayout)
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

	nowParsed, err := time.Parse(dateLayout, now)
	if err != nil {
		return "", err
	}

	dateTime, err := time.Parse(dateLayout, dstart)
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
	case "w":
		return nextWeekdayDate(nowParsed, dateTime, repeatRules)
	case "m":
		return nextMonthDate(nowParsed, dateTime, repeatRules)
	default:
		return "", errors.New("invalid rule for next date: "+repeatRules[0])
	}	

	return resultDate.Format(dateLayout), nil
}

func nextWeekdayDate(currentTime, startDate time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("weekday rule requires days specification")
	}

	weekdays := [8]bool{}
	
	for _, dayStr := range strings.Split(parts[1], ",") {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", fmt.Errorf("invalid weekday: %s", dayStr)
		}
		if day < 1 || day > 7 {
			return "", fmt.Errorf("weekday must be between 1 and 7, got: %d", day)
		}
		weekdays[day] = true
	}

	searchDate := currentTime.AddDate(0, 0, 1)
	
	for i := 0; i < 730; i++ {
		currentWeekday := int(searchDate.Weekday())
		if currentWeekday == 0 {
			currentWeekday = 7
		}

		if weekdays[currentWeekday] && searchDate.After(startDate) {
			return searchDate.Format(dateLayout), nil
		}
		searchDate = searchDate.AddDate(0, 0, 1)
	}

	return "", errors.New("no suitable date found within 2 years")
}

func nextMonthDate(currentTime, startDate time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("month rule requires days specification")
	}

	normalDays := [32]bool{}    
	specialDays := [3]bool{}    
	months := [13]bool{}        

	for _, dayStr := range strings.Split(parts[1], ",") {
		day, err := strconv.Atoi(strings.TrimSpace(dayStr))
		if err != nil {
			return "", fmt.Errorf("invalid day: %s", dayStr)
		}
		
		if day == -1 {
			specialDays[1] = true
		} else if day == -2 {
			specialDays[0] = true
		} else if day >= 1 && day <= 31 {
			normalDays[day] = true
		} else {
			return "", fmt.Errorf("day must be between 1-31, -1, or -2, got: %d", day)
		}
	}

	if len(parts) >= 3 {
		for _, monthStr := range strings.Split(parts[2], ",") {
			month, err := strconv.Atoi(strings.TrimSpace(monthStr))
			if err != nil {
				return "", fmt.Errorf("invalid month: %s", monthStr)
			}
			if month < 1 || month > 12 {
				return "", fmt.Errorf("month must be between 1 and 12, got: %d", month)
			}
			months[month] = true
		}
	} else {
		for i := 1; i <= 12; i++ {
			months[i] = true
		}
	}

	searchDate := currentTime.AddDate(0, 0, 1)
	
	for i := 0; i < 730; i++ {
		currentYear := searchDate.Year()
		currentMonth := int(searchDate.Month())
		currentDay := searchDate.Day()
		daysInCurrentMonth := daysInMonth(currentYear, currentMonth)

		if !months[currentMonth] {
			searchDate = searchDate.AddDate(0, 0, 1)
			continue
		}

		found := false
		
		if currentDay <= 31 && normalDays[currentDay] && currentDay <= daysInCurrentMonth {
			found = true
		}
		
		if specialDays[1] && currentDay == daysInCurrentMonth {
			found = true
		}
		if specialDays[0] && currentDay == daysInCurrentMonth-1 {
			found = true
		}

		if found && searchDate.After(startDate) {
			return searchDate.Format(dateLayout), nil
		}

		searchDate = searchDate.AddDate(0, 0, 1)
	}

	return "", errors.New("no suitable date found within 2 years")
}

func daysInMonth(year, month int) int {
	return time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()
}