package api

import (
	"fmt"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/makmanu/go_final_project/pkg/db"
)

var task *db.Task

type taskError struct {
    Error      string `json:"error,omitempty"`
}
var jsonError taskError

const dateLayout = "20060102"

func checkDate(task *db.Task) error {
	now := time.Now()
	if len(task.Date) == 0 {
            task.Date = now.Format(dateLayout)
	}
	if task.Date == now.Format(dateLayout) {
			return nil
	}
	t, err := time.Parse(dateLayout, task.Date)
	if err != nil {
		return err
	}
	if afterNow(now, t) {
        if len(task.Repeat) == 0 {
            task.Date = now.Format(dateLayout)
        } else {
            task.Date, err = nextDate(now.Format(dateLayout), task.Date, task.Repeat)		
			if err != nil {
				return err
			}
        }
    }
	return nil
}

func writeJson(w http.ResponseWriter, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		writeJson(w, jsonError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	log.Println("response with:", data)
	w.Write(resp)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	log.Println("start to handle /api/task post")

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	log.Print("\nget task:", task.Title, "\nwith comment: ", task.Comment, "\ndate: ", task.Date, "\nrepeat rule: ", task.Repeat)
	if len(task.Title) == 0 {
		jsonError.Error = "no Title"
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	if err := checkDate(task); err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	Id, err := db.AddTask(task)
	task.ID = fmt.Sprint(Id)
	if err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	writeJson(w, task)
	w.WriteHeader(http.StatusOK)
}