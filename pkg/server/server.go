package server

import (
	"log"
	"net/http"
)

func Start_server(port string) error {
	webDir := "./web"
	log.Println("Запуск сервера")
	defer log.Println("Остановка сервера")

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	err := http.ListenAndServe(port, nil)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}