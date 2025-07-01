package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fuckingUnicorn/final_project/pkg/constants"
	"github.com/fuckingUnicorn/final_project/pkg/db"
	"github.com/fuckingUnicorn/final_project/pkg/entities"
)

func getTaskListHandler(w http.ResponseWriter, r *http.Request) {

	var filter entities.Filter
	filter.Limit = constants.DefaultLimit

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		offset, err := strconv.ParseInt(offsetParam, 10, 64)
		if err != nil {
			err = errors.New("failed to parse 'offset' parameter")
			writeJSON(w, http.StatusBadRequest, err)
			return
		}
		filter.Offset = offset
	}

	taskList, err := db.GetTaskList(&filter)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err)
		return
	}

	if taskList == nil {
		writeJSON(w, http.StatusOK, []entities.Task{})
	} else {
		writeJSON(w, http.StatusOK, taskList)
	}
}
