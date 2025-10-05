package api

import ("net/http"
		"github.com/makmanu/go_final_project/pkg/db")


type TasksResp struct {
    Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request){
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