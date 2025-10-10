package server

import (
	"log"
	"net/http"
	"os"

	"github.com/makmanu/go_final_project/pkg/api"
)

func Start_server(port string) error {
	Password := os.Getenv("TODO_PASSWORD")
    if Password != "" {
        log.Println("Authentication is enabled")
    } else {
        log.Println("Authentication is disabled - TODO_PASSWORD not set")
    }
    
    http.HandleFunc("/api/signin", api.SignInHandler)

	webDir := "./web"
	log.Println("Запуск сервера")
	defer log.Println("Остановка сервера")

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", api.HandleQuery)
	http.HandleFunc("/api/task", api.AuthMiddleware(api.TaskHandler))
	http.HandleFunc("/api/tasks", api.AuthMiddleware(api.TasksHandler))
	http.HandleFunc("/api/task/done", api.AuthMiddleware(api.TaskDone))

	err := http.ListenAndServe(port, nil)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}
	
