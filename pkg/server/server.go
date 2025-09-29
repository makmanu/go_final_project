package server

import (
	"log"
	"net/http"

	"github.com/makmanu/go_final_project/pkg/api"
)

func Start_server(port string) error {
	webDir := "./web"
	log.Println("Запуск сервера")
	defer log.Println("Остановка сервера")

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", api.HandleQuery)
	http.HandleFunc("/api/task", api.TaskHandler)

	err := http.ListenAndServe(port, nil)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}
	
