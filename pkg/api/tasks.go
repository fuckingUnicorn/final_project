package api

import (
	"log"
	"net/http"

	"go1f/pkg/db"
)

// TasksHandler обрабатывает запросы к /api/tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	log.Printf("TasksHandler: получен %s запрос к %s", r.Method, r.URL.String())

	// Получаем параметр поиска из запроса
	search := r.URL.Query().Get("search")
	log.Printf("TasksHandler: параметр поиска: %s", search)

	// Получаем список задач из базы данных
	tasks, err := db.Tasks(search)
	if err != nil {
		log.Printf("TasksHandler: ошибка при получении задач: %v", err)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	log.Printf("TasksHandler: получено %d задач", len(tasks))

	// Преобразуем список задач в формат, понятный тестам
	result := make([]map[string]string, len(tasks))
	for i, task := range tasks {
		result[i] = map[string]string{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
	}

	// Оборачиваем результат в map с ключом "tasks", как ожидает тест
	response := map[string][]map[string]string{
		"tasks": result,
	}

	log.Printf("TasksHandler: возвращаем ответ с %d задачами", len(result))
	writeJSON(w, response)
}
