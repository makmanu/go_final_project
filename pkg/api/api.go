package api

import "net/http"


func TaskHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
        addTaskHandler(w, r)
    case http.MethodGet:
        getTaskById(w, r)
    }
}