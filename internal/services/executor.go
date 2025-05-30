package services

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"ikurotime/code-engine/internal/models"
)

type Executor struct {
	semaphore chan struct{}
	timeout   time.Duration
}

func NewExecutor(maxConcurrent int, timeout time.Duration) *Executor {
	return &Executor{
		semaphore: make(chan struct{}, maxConcurrent),
		timeout:   timeout,
	}
}

func (e *Executor) Execute(req models.ExecuteRequest) (string, error) {
	// Acquire semaphore
	e.semaphore <- struct{}{}
	defer func() { <-e.semaphore }()

	// Create temp directory
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("codeexec_%d_%d", time.Now().UnixNano(), rand.Int63()))
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Write code to file
	fileName := filepath.Join(tempDir, fmt.Sprintf("code.%s", models.LanguageToExtension[req.Language]))
	if err := os.WriteFile(fileName, []byte(req.Code), 0644); err != nil {
		return "", fmt.Errorf("failed to write code file: %w", err)
	}

	// Execute with timeout
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, req.Language, fileName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("execution failed: %w", err)
	}

	return string(output), nil
}
