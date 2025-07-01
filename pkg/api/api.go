package api

import (
	"github.com/go-chi/chi/v5"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", getNextDayHandler)
	r.Post("/api/task", addTaskHandler)
	r.Get("/api/tasks", getTaskListHandler)
	r.Get("/api/task", getTaskByIdHandler)
	r.Put("/api/task", updateTaskHandler)
	r.Delete("/api/task", deleteTaskHandler)
	r.Post("/api/task/done", CompleteTaskHandler)
}
