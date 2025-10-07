package api

import (
	"log"
	"net/http"

	"github.com/makmanu/go_final_project/pkg/db"
)


type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request){
	log.Println("start to handle /api/tasks")
    tasks, err := db.Tasks(50)
    if err != nil {
        jsonError.Error = err.Error()
		writeJson(w, jsonError)
        return
    }
    writeJson(w, TasksResp{
        Tasks: tasks,
    })
}