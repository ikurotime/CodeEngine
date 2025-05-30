package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

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

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("%s %s", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *Handler) Execute(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("%s %s", r.Method, r.URL.Path)

	if r.Method != "POST" {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("Error reading request body: %s", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to read request body")
		return
	}

	var request models.ExecuteRequest
	if err := json.Unmarshal(body, &request); err != nil {
		h.logger.Printf("Error unmarshalling request: %s", err)
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	h.logger.Printf("Request: %+v", request)

	output, err := h.executor.Execute(request)
	if err != nil {
		h.logger.Printf("Error executing code: %s", err)
		h.writeErrorResponse(w, http.StatusInternalServerError, "Execution failed")
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
