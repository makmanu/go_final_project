package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/makmanu/go_final_project/pkg/db"
)


func updateTask(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	log.Println("start to handle /api/task put")

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
	err = db.UpdateTask(task)
	if err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	writeJson(w, task)
	w.WriteHeader(http.StatusOK)
}