package main

import (
	"log"
	"os"

	"github.com/makmanu/go_final_project/pkg/db"
	"github.com/makmanu/go_final_project/pkg/server"
)

func main() {
	envDBPath := os.Getenv("TODO_DBFILE")
	if envDBPath == "" {
		envDBPath = "pkg/db/scheduler.db"
	}
	err := db.Init(envDBPath)
	if err != nil{
		log.Fatalln("Database initialization err:", err)
	}
	envPort := os.Getenv("TODO_PORT")
	if envPort== "" {
		envPort = "7540"
	}
	server.Start_server(":"+envPort)
}