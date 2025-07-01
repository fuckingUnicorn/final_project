package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fuckingUnicorn/final_project/pkg/db"
	"github.com/fuckingUnicorn/final_project/pkg/entities"
	"github.com/fuckingUnicorn/final_project/pkg/internal"
)

// updateTaskHandler - HTTP-обработчик для обновления существующей задачи
// Метод: PUT

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *entities.Task
	var buf bytes.Buffer

	// 1. Чтение тела запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err = errors.New("Error reading body: " + err.Error())
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// 2. Парсинг JSON
	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		err = errors.New("Error unmarshalling body: " + err.Error())
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// 3. Валидация обязательных полей
	if task.Title == "" {
		err = errors.New("title is required")
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// 4. Проверка даты
	if err = internal.CheckDate(task); err != nil {
		err = errors.New("checkDate failed: " + err.Error())
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// 5. Обновление в БД
	err = db.UpdateTask(task)
	if err != nil {
		err = errors.New("Error adding task: " + err.Error())
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// 6. Успешный ответ
	writeJSON(w, http.StatusOK, entities.EmptyResponse{})
}
