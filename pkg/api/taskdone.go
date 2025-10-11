package api

import (
	"log"
	"net/http"
	"time"

	"github.com/makmanu/go_final_project/pkg/db"
)


func TaskDone(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		jsonError.Error = "no id"
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	id := r.URL.Query().Get("id")
	currentTask, err := db.GetTask(id)
	if err != nil{
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	if len(currentTask.Repeat) == 0{
		err = db.DeleteTask(currentTask.ID)
		if err != nil {
			jsonError.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			writeJson(w, jsonError)
			return
		}
		log.Println("deleted task with id ", id)
		jsonError.Error = ""
		w.WriteHeader(http.StatusOK)
		writeJson(w, jsonError)
		return
	}
	currentTask.Date, err = nextDate(time.Now().Format(dateLayout), currentTask.Date, currentTask.Repeat)
	if err != nil{
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	err = db.UpdateTask(currentTask)
	if err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	jsonError.Error = ""
	w.WriteHeader(http.StatusOK)
	writeJson(w, jsonError)
}