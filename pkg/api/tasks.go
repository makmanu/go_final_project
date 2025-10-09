package api

import (
	"log"
	"net/http"

	"github.com/makmanu/go_final_project/pkg/db"
)


type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("start to handle /api/tasks")
	
	// Получаем параметр search из query string
	search := r.URL.Query().Get("search")
	
	var tasks []*db.Task
	var err error
	
	if search != "" {
		tasks, err = db.SearchTasks(search, 50)
	} else {
		tasks, err = db.Tasks(50)
	}
	
	if err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}