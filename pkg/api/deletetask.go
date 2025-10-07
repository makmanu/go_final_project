package api

import (
	"log"
	"net/http"

	"github.com/makmanu/go_final_project/pkg/db"
)


func deleteTask(w http.ResponseWriter, r *http.Request) {
	log.Println("start to handle /api/task delete")
	if !r.URL.Query().Has("id") {
		jsonError.Error = "no id"
		writeJson(w, jsonError)
		return
	}
	id := r.URL.Query().Get("id")
	err := db.DeleteTask(id)
	if err != nil{
		jsonError.Error = err.Error()
		writeJson(w, jsonError)
		return
	}
	writeJson(w, jsonError)
}