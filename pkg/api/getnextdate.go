package api

import (
	"net/http"
	"time"

	"github.com/fuckingUnicorn/final_project/pkg/constants"
	"github.com/fuckingUnicorn/final_project/pkg/internal"
)

// getNextDayHandler - HTTP-обработчик для расчета следующей даты выполнения
// Метод: GET

func getNextDayHandler(w http.ResponseWriter, r *http.Request) {
	var nowDate time.Time

	// 1. Обработка параметра now
	nowParam := r.URL.Query().Get("now")
	if nowParam != "" {
		var err error
		nowDate, err = time.Parse(constants.DateFormat, nowParam)
		if err != nil {
			nowDate = time.Now() // Используем текущую дату при ошибке парсинга
		}
	} else {
		nowDate = time.Now() // Используем текущую дату по умолчанию
	}

	// 2. Получение основных параметров
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// 3. Расчет следующей даты
	nextDate, err := internal.NextDate(nowDate, date, repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	// 4. Отправка результата
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
