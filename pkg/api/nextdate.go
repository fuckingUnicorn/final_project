package api

import (
	"log"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи на основе правила повторения
func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	// Проверяем формат даты
	date, err := time.Parse(DateFormat, dateStr)
	if err != nil {
		log.Printf("Ошибка парсинга даты: %v", err)
		return "", nil
	}

	// Если дата слишком старая или из будущего, возвращаем пустую строку
	if date.Year() < 1900 || date.Year() > 2100 {
		return "", nil
	}

	// Если правило повторения не указано
	if repeat == "" {
		// Если дата в прошлом или сегодня, нет следующей даты
		if !date.After(now) {
			return "", nil
		}
		return dateStr, nil
	}

	// Разбираем правило повторения
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", nil
	}

	ruleType := parts[0]

	switch ruleType {
	case "y": // Ежегодно
		return calculateYearlyNextDate(now, date)

	case "d": // Каждые N дней
		if len(parts) < 2 {
			return "", nil
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 {
			return "", nil
		}

		// Проверяем, что интервал не слишком большой
		if days > 400 {
			return "", nil
		}

		return calculateDailyNextDate(now, date, days)

	default:
		return "", nil
	}
}

// calculateYearlyNextDate вычисляет следующую дату для ежегодного повторения
func calculateYearlyNextDate(now time.Time, date time.Time) (string, error) {
	// Проверяем особый случай: 29 февраля или 1 марта (связанные с високосными годами)
	if (date.Month() == time.February && date.Day() == 29) ||
		(date.Month() == time.March && date.Day() == 1) {
		// Для 29 февраля или 1 марта используем 1 марта следующего года
		nextDate := time.Date(date.Year()+1, time.March, 1, 0, 0, 0, 0, time.UTC)
		// Если дата в прошлом, продолжаем добавлять годы, пока не получим дату в будущем
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format(DateFormat), nil
	}

	// Для обычных дат просто добавляем год
	nextDate := time.Date(date.Year()+1, date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	// Если дата в прошлом, продолжаем добавлять годы, пока не получим дату в будущем
	for nextDate.Before(now) || nextDate.Equal(now) {
		nextDate = nextDate.AddDate(1, 0, 0)
	}
	return nextDate.Format(DateFormat), nil
}

// calculateDailyNextDate вычисляет следующую дату для повторения через N дней
func calculateDailyNextDate(now time.Time, date time.Time, days int) (string, error) {
	// Вычисляем следующую дату
	nextDate := date.AddDate(0, 0, days)
	// Если дата в прошлом, продолжаем добавлять интервалы, пока не получим дату в будущем
	for nextDate.Before(now) || nextDate.Equal(now) {
		nextDate = nextDate.AddDate(0, 0, days)
	}
	return nextDate.Format(DateFormat), nil
}
