package api

import (
	"log"
	"net/http"
	"time"
)

// NextDateHandler обрабатывает GET-запросы к /api/nextdate
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	now := time.Now()
	if nowStr != "" {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeatStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(nextDate))
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	// Передаем все запросы в TaskHandler
	log.Printf("Запрос к /api/task: %s %s", r.Method, r.URL.String())
	TaskHandler(w, r)
}
