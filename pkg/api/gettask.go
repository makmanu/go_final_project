package api

import (
	"net/http"

	"github.com/makmanu/go_final_project/pkg/db"
)

func getTaskById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		jsonError.Error = "missing ID"
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		jsonError.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		writeJson(w, jsonError)
		return
	}
	writeJson(w, task)
	w.WriteHeader(http.StatusOK)
}