package http

import (
	"fmt"
	"net/http"
)

func runningTasks() int {
	if taskRunner == nil {
		return 0
	}
	return taskRunner.ActiveWorkers()
}

func totalTasks() int {
	if taskRunner == nil {
		return 0
	}
	return taskRunner.TotalWorkers()
}

func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(fmt.Sprintf(`{"running_tasks": %d, "total_tasks": %d}`, runningTasks(), totalTasks())))
	if err != nil {
		logger.Error(err, " Failed to write healthcheck")
	}
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"ready": true}`))
	if err != nil {
		logger.Error(err, " Failed to write ready")
	}
}

func handleLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"live": true}`))
	if err != nil {
		logger.Error(err, " Failed to write live")
	}
}
