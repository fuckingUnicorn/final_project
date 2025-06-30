package api

import "net/http"

// Init регистрирует все API-обработчики
func Init() {
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", TaskDoneHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)
}
