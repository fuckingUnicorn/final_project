package api

import (
	"log"
	"net/http"
	"time"

	"go1f/pkg/db"
)

// TaskDoneHandler обрабатывает запросы к /api/task/done для отметки задачи как выполненной
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	
	log.Printf("TaskDoneHandler: получен %s запрос к %s", r.Method, r.URL.String())
	
	// Проверяем, что метод запроса - POST
	if r.Method != http.MethodPost {
		log.Printf("TaskDoneHandler: метод %s не поддерживается", r.Method)
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"})
		return
	}
	
	// Получаем ID задачи из параметров запроса
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("TaskDoneHandler: не указан ID задачи")
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	
	log.Printf("TaskDoneHandler: запрос на отметку выполнения задачи с ID: %s", id)
	
	// Получаем задачу по ID
	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("TaskDoneHandler: ошибка при получении задачи: %v", err)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	// Проверяем, что задача не nil
	if task == nil {
		log.Printf("TaskDoneHandler: задача с ID %s не найдена", id)
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}
	
	log.Printf("TaskDoneHandler: получена задача: ID=%s, Date=%s, Title=%s, Comment=%s, Repeat=%s", 
		task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	
	// Если задача не имеет правила повторения, удаляем её
	if task.Repeat == "" {
		log.Printf("TaskDoneHandler: задача без повторения, удаляем")
		if err := db.DeleteTask(id); err != nil {
			log.Printf("TaskDoneHandler: ошибка при удалении задачи: %v", err)
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
		log.Printf("TaskDoneHandler: задача успешно удалена")
		writeJSON(w, map[string]string{})
		return
	}
	
	// Если задача периодическая, вычисляем следующую дату
	log.Printf("TaskDoneHandler: задача с повторением, вычисляем следующую дату")
	now := time.Now() // Текущее время
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		log.Printf("TaskDoneHandler: ошибка при вычислении следующей даты: %v", err)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	log.Printf("TaskDoneHandler: следующая дата: %s", next)
	
	// Обновляем дату задачи
	if err := db.UpdateDate(next, id); err != nil {
		log.Printf("TaskDoneHandler: ошибка при обновлении даты задачи: %v", err)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}
	
	log.Printf("TaskDoneHandler: дата задачи успешно обновлена")
	writeJSON(w, map[string]string{})
}
