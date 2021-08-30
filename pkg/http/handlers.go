package http

import (
	"fmt"
	"net/http"
)

func (h *Server) runningTasks() int {
	if h.taskRunner == nil {
		return 0
	}
	return h.taskRunner.ActiveWorkers()
}

func (h *Server) totalTasks() int {
	if h.taskRunner == nil {
		return 0
	}
	return h.taskRunner.TotalWorkers()
}

func (h *Server) handleHealthcheck(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(fmt.Sprintf(`{"running_tasks": %d, "total_tasks": %d}`, h.runningTasks(), h.totalTasks())))
	if err != nil {
		h.logger.Error(err, " Failed to write healthcheck")
	}
}

func (h *Server) handleReady(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"ready": true}`))
	if err != nil {
		h.logger.Error(err, " Failed to write ready")
	}
}

func (h *Server) handleLive(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, err := w.Write([]byte(`{"live": true}`))
	if err != nil {
		h.logger.Error(err, " Failed to write live")
	}
}
