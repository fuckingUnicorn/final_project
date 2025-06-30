package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go1f/pkg/db"
)

// writeJSON сериализует данные в JSON и отправляет их в ответ
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка сериализации JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// checkDate проверяет и корректирует дату задачи в соответствии с правилами повторения
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана или указана как 'today', используем текущую дату
	if task.Date == "" || task.Date == "today" {
		task.Date = now.Format(DateFormat)
	}

	// Проверяем формат даты
	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return err
	}

	// Если дата меньше текущей и правило повторения не указано, используем текущую дату
	if t.Before(now.Truncate(24*time.Hour)) && task.Repeat == "" {
		task.Date = now.Format(DateFormat)
		return nil
	}

	// Если дата в прошлом и есть правило повторения, вычисляем следующую дату
	if t.Before(now.Truncate(24*time.Hour)) && task.Repeat != "" {
		// Используем общую функцию NextDate для вычисления следующей даты
		nextDateStr, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}

		// Если NextDate вернула пустую строку, это означает, что нет следующей даты
		if nextDateStr == "" {
			return fmt.Errorf("дата не может быть меньше сегодняшней")
		}

		task.Date = nextDateStr
		log.Printf("Вычислена следующая дата: %s", task.Date)
	}

	return nil
}
