package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"ikurotime/code-engine/internal/models"
	"ikurotime/code-engine/internal/services"
)

type Handler struct {
	executor *services.Executor
	logger   *log.Logger
}

func NewHandler(executor *services.Executor, logger *log.Logger) *Handler {
	return &Handler{
		executor: executor,
		logger:   logger,
	}
}

func (h *Handler) Execute(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != "POST" {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Check if executor is shutting down
	if h.executor.IsShutdown() {
		h.writeErrorResponse(w, http.StatusServiceUnavailable, "Service is shutting down")
		return
	}

	//formdata
	code := r.FormValue("code")
	language := r.FormValue("language")

	if code == "" || language == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Code and language are required")
		return
	}

	request := models.ExecuteRequest{
		Code:     code,
		Language: language,
	}

	h.logger.Printf("Request: %+v", request)

	output, err := h.executor.Execute(request)
	if err != nil {
		h.logger.Printf("Error executing code: %s, output: %s", err, output)

		// Check if error is due to shutdown
		if strings.Contains(err.Error(), "shutting down") {
			h.writeErrorResponse(w, http.StatusServiceUnavailable, "Service is shutting down")
			return
		}

		// Return execution error with output
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(models.ExecuteResponse{Output: output})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ExecuteResponse{Output: output})
}

func (h *Handler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{Error: message})
}
