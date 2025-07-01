package db

import (
	"database/sql"
	"fmt"
	"log"
)

// Task представляет задачу в планировщике
type Task struct {
	ID      string
	Date    string
	Title   string
	Comment string
	Repeat  string
}

// AddTask добавляет новую задачу в базу данных
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// GetTask возвращает задачу по её идентификатору
func GetTask(id string) (*Task, error) {
	log.Printf("db.GetTask: запрос задачи с ID: %s", id)
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	log.Printf("db.GetTask: выполняем запрос: %s с параметром %s", query, id)

	row := DB.QueryRow(query, id)

	// Создаем временные переменные для сканирования
	var taskID int64
	var date, title, comment, repeat string

	// Сканируем значения в переменные
	err := row.Scan(&taskID, &date, &title, &comment, &repeat)
	if err == sql.ErrNoRows {
		log.Printf("db.GetTask: задача не найдена")
		return nil, fmt.Errorf("задача не найдена")
	} else if err != nil {
		log.Printf("db.GetTask: ошибка при получении задачи: %v", err)
		return nil, fmt.Errorf("ошибка при получении задачи: %v", err)
	}

	log.Printf("db.GetTask: получены данные из БД: ID=%d, Date=%s, Title=%s, Comment=%s, Repeat=%s",
		taskID, date, title, comment, repeat)

	// Создаем и возвращаем задачу с преобразованным ID
	task := &Task{
		ID:      fmt.Sprint(taskID),
		Date:    date,
		Title:   title,
		Comment: comment,
		Repeat:  repeat,
	}

	log.Printf("db.GetTask: возвращаем задачу: %+v", task)
	return task, nil
}

// UpdateTask обновляет задачу в базе данных
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// DeleteTask удаляет задачу из базы данных
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id=?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// UpdateDate обновляет только дату задачи в базе данных
func UpdateDate(next string, id string) error {
	log.Printf("db.UpdateDate: обновление даты задачи %s на %s", id, next)
	query := `UPDATE scheduler SET date=? WHERE id=?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		log.Printf("db.UpdateDate: ошибка при обновлении даты: %v", err)
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		log.Printf("db.UpdateDate: ошибка при получении числа затронутых строк: %v", err)
		return err
	}
	if n == 0 {
		log.Printf("db.UpdateDate: задача не найдена")
		return fmt.Errorf("задача не найдена")
	}
	log.Printf("db.UpdateDate: дата задачи успешно обновлена")
	return nil
}

// Tasks возвращает список задач из базы данных
func Tasks(search string) ([]Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler`
	var rows *sql.Rows
	var err error
	if search != "" {
		query += ` WHERE title LIKE ? OR comment LIKE ? OR date LIKE ?`
		search = "%" + search + "%"
		rows, err = DB.Query(query, search, search, search)
	} else {
		rows, err = DB.Query(query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var id int64
		var date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, Task{
			ID:      fmt.Sprint(id),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}
	return tasks, rows.Err()
}
