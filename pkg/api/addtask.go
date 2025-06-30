package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"go1f/pkg/db"
)

// Используем константу DateFormat из файла constants.go

// AddTaskHandler обрабатывает POST-запросы к /api/task для добавления новой задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Проверяем, что метод POST
	if r.Method != http.MethodPost {
		writeJSON(w, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	// Обрабатываем POST-запрос для добавления новой задачи
	addTask(w, r)
}

// addTask обрабатывает добавление новой задачи
func addTask(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос на добавление задачи получен")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Ошибка чтения тела запроса: %v", err)
		writeJSON(w, map[string]interface{}{"error": "Неверный формат данных"})
		return
	}
	log.Printf("Тело запроса: %s", string(body))

	var task db.Task

	if err := json.Unmarshal(body, &task); err != nil {
		log.Printf("Ошибка парсинга JSON: %v", err)
		writeJSON(w, map[string]interface{}{"error": "Неверный формат данных"})
		return
	}

	log.Printf("Получена задача: %+v", task)

	// Проверяем обязательное поле title
	if task.Title == "" {
		writeJSON(w, map[string]interface{}{"error": "Не указан заголовок задачи"})
		return
	}

	// Обрабатываем дату и правило повторения
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]interface{}{"error": err.Error()})
		return
	}

	// Добавляем задачу в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		log.Printf("Ошибка при добавлении задачи: %v", err)
		writeJSON(w, map[string]interface{}{"error": "Ошибка при добавлении задачи"})
		return
	}

	log.Printf("Задача добавлена с id: %d", id)
	writeJSON(w, map[string]interface{}{
		"id": strconv.FormatInt(id, 10),
	})
}

// Функции checkDate и writeJSON перенесены в utils.go
