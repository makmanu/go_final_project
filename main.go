package main

import (
	"fmt"
	"log"
	"time"

	"github.com/makmanu/go_final_project/pkg/db"
	"github.com/makmanu/go_final_project/pkg/server"
)

func main() {
	err := db.Init("pkg/db/scheduler.db")
	if err != nil{
		log.Fatalln("Database initialization err:", err)
	}
	server.Start_server(":7540")
}