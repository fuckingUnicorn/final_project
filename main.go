package main

import (
	"log"
	"os"

	"go1f/pkg/db"
	"go1f/pkg/server"
)

func main() {
	// Определяем путь к базе данных: либо из переменной, либо по умолчанию
	dbFile := "scheduler.db"
	if envDb := os.Getenv("TODO_DBFILE"); envDb != "" {
		dbFile = envDb
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("db init error: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
