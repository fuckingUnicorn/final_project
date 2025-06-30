package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"go1f/pkg/db"
)

// TaskHandler обрабатывает GET, PUT и DELETE запросы к /api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	log.Printf("TaskHandler: получен %s запрос к %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		// Получаем ID задачи из параметров запроса
		id := r.URL.Query().Get("id")
		if id == "" {
			log.Printf("TaskHandler: не указан ID задачи")
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}

		// Получаем задачу по ID
		log.Printf("TaskHandler: запрос на получение задачи с ID: %s", id)

		// Проверяем, что ID - это число
		_, parseErr := strconv.ParseInt(id, 10, 64)
		if parseErr != nil {
			log.Printf("TaskHandler: неверный формат ID: %s, ошибка: %v", id, parseErr)
			writeJSON(w, map[string]string{"error": "Неверный формат ID"})
			return
		}

		task, err := db.GetTask(id)
		if err != nil {
			log.Printf("TaskHandler: ошибка при получении задачи: %v", err)
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		log.Printf("TaskHandler: получена задача из БД: %+v", task)

		// Проверяем, что задача не nil
		if task == nil {
			log.Printf("TaskHandler: ошибка: задача с ID %s не найдена (но GetTask не вернул ошибку)", id)
			writeJSON(w, map[string]string{"error": "Задача не найдена"})
			return
		}

		log.Printf("TaskHandler: полученная задача: ID=%s, Date=%s, Title=%s, Comment=%s, Repeat=%s",
			task.ID, task.Date, task.Title, task.Comment, task.Repeat)

		// Возвращаем задачу в формате map[string]string для совместимости с тестами
		response := map[string]string{
			"id":      task.ID,
			"date":    task.Date,
			"title":   task.Title,
			"comment": task.Comment,
			"repeat":  task.Repeat,
		}
		log.Printf("TaskHandler: возвращаем ответ: %+v", response)
		writeJSON(w, response)

	case http.MethodPut:
		// Декодируем JSON из тела запроса
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSON(w, map[string]string{"error": "Неверный формат данных"})
			return
		}

		// Проверяем обязательное поле title
		if task.Title == "" {
			writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
			return
		}

		// Обрабатываем дату и правило повторения
		if err := checkDate(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		// Обновляем задачу в базе данных
		if err := db.UpdateTask(&task); err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, map[string]string{"id": task.ID})

	case http.MethodDelete:
		// Получаем ID задачи из параметров запроса
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
			return
		}

		// Удаляем задачу из базы данных
		if err := db.DeleteTask(id); err != nil {
			log.Printf("TaskHandler: ошибка при удалении задачи: %v", err)
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		log.Printf("TaskHandler: задача с ID %s успешно удалена", id)
		writeJSON(w, map[string]string{})

	case http.MethodPost:
		// Делегируем обработку POST-запроса функции AddTaskHandler
		AddTaskHandler(w, r)

	default:
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"})
	}
}
