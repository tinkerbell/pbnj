package http

import (
	"fmt"
	"net/http"
)

func (h *HTTPServer) runningTasks() int {
	if h.taskRunner == nil {
		return 0
	}
	return h.taskRunner.ActiveWorkers()
}

func (h *HTTPServer) totalTasks() int {
	if h.taskRunner == nil {
		return 0
	}
	return h.taskRunner.TotalWorkers()
}

func (h *HTTPServer) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(fmt.Sprintf(`{"running_tasks": %d, "total_tasks": %d}`, h.runningTasks(), h.totalTasks())))
	if err != nil {
		h.logger.Error(err, " Failed to write healthcheck")
	}
}

func (h *HTTPServer) handleReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"ready": true}`))
	if err != nil {
		h.logger.Error(err, " Failed to write ready")
	}
}

func (h *HTTPServer) handleLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"live": true}`))
	if err != nil {
		h.logger.Error(err, " Failed to write live")
	}
}
