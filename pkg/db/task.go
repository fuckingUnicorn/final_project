package db

import (
	"errors"
	"strconv"

	"github.com/fuckingUnicorn/final_project/pkg/constants"
	"github.com/fuckingUnicorn/final_project/pkg/entities"
)

func AddTask(task *entities.Task) (int64, error) {
	var id int64
	res, err := db.Exec(constants.QueryAddTask, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func GetTaskList(filter *entities.Filter) ([]entities.Task, error) {

	rows, err := db.Query(constants.QueryGetTaskList, constants.SortASC, filter.Limit, filter.Offset)
	if err != nil {
		err = errors.New("failed to db.GetTaskList: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	var tasks []entities.Task

	for rows.Next() {
		var task entities.Task
		err = rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err = rows.Err(); err != nil {
			err = errors.New("failed to Scan rows in Task entity: " + err.Error())
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func GetTaskById(id string) (*entities.Task, error) {
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = errors.New("failed to parse task id")
		return nil, err
	}

	var task entities.Task
	err = db.QueryRow(constants.QueryGetTaskById, taskId).Scan(&taskId, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		err = errors.New("failed to db.GetTaskBy: " + err.Error())
		return nil, err
	}

	task.ID = strconv.Itoa(int(taskId))

	return &task, nil
}

func UpdateTask(task *entities.Task) error {
	taskId, err := strconv.ParseInt(task.ID, 10, 64)
	if err != nil {
		err = errors.New("failed to parse task id")
		return err
	}

	res, err := db.Exec(constants.QueryUpdateTask, taskId, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	taskId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		err = errors.New("failed to parse task id")
		return err
	}
	res, err := db.Exec(constants.QueryDeleteTask, taskId)
	if err != nil {
		err = errors.New("failed to delete task: " + err.Error())
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		err = errors.New("failed to count affected rows: " + err.Error())
		return err
	}
	if count == 0 {
		err = errors.New(`incorrect id for deleting task: ` + id)
		return err
	}

	return nil
}
