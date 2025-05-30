package services

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"ikurotime/code-engine/internal/models"
)

type Executor struct {
	pools   map[string]*ContainerPool
	timeout time.Duration
	logger  *log.Logger
}

func NewExecutor(maxConcurrent int, timeout time.Duration, logger *log.Logger) *Executor {
	executor := &Executor{
		pools:   make(map[string]*ContainerPool),
		timeout: timeout,
		logger:  logger,
	}

	for language, image := range models.LanguageToSandbox {
		pool := &ContainerPool{
			containers: make(chan string, maxConcurrent),
			language:   language,
			image:      image,
			maxSize:    maxConcurrent,
			logger:     logger,
		}
		executor.pools[language] = pool

		executor.logger.Printf("Initializing container pool for %s (image: %s, size: %d)", language, image, maxConcurrent)

		go executor.initializePool(pool)
	}

	return executor
}

func (e *Executor) Execute(req models.ExecuteRequest) (string, error) {
	pool, exists := e.pools[req.Language]
	if !exists {
		return "", fmt.Errorf("unsupported language: %s", req.Language)
	}

	e.logger.Printf("Executing %s code, waiting for container from pool...", req.Language)

	containerID, err := getContainerFromPool(pool, req.Language)
	if err != nil {
		return "", fmt.Errorf("failed to get container from pool: %w", err)
	}

	fileName, err := e.createTempFiles(req)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	if err := e.copyCodeToContainer(containerID, fileName, req.Language); err != nil {
		return "", fmt.Errorf("failed to copy code to container: %w", err)
	}

	e.logger.Printf("Executing %s code in container %s", req.Language, containerID[:12])

	output, err := e.executeCodeInContainer(containerID, req.Language, fileName)

	if err != nil {
		e.logger.Printf("Code execution failed in container %s: %v", containerID[:12], err)
		return string(output), fmt.Errorf("execution failed: %w", err)
	}

	e.logger.Printf("Code execution completed successfully in container %s", containerID[:12])
	return string(output), nil
}

func getContainerFromPool(pool *ContainerPool, language string) (string, error) {
	var containerID string
	select {
	case containerID = <-pool.containers:
		pool.logger.Printf("Acquired container %s for %s execution", containerID[:12], language)
		defer func() {
			pool.containers <- containerID
			pool.logger.Printf("Returned container %s to %s pool", containerID[:12], language)
		}()
	case <-time.After(5 * time.Second):
		return "", fmt.Errorf("no containers available, pool exhausted")
	}
	return containerID, nil
}

func (e *Executor) createTempFiles(req models.ExecuteRequest) (string, error) {
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("codeexec_%d_%d", time.Now().UnixNano(), rand.Int63()))
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	fileName := filepath.Join(tempDir, fmt.Sprintf("script.%s", models.LanguageToExtension[req.Language]))
	if err := os.WriteFile(fileName, []byte(req.Code), 0644); err != nil {
		return "", fmt.Errorf("failed to write code file: %w", err)
	}

	return fileName, nil
}

func (e *Executor) copyCodeToContainer(containerID string, fileName string, language string) error {
	copyCmd := exec.Command("docker", "cp", fileName, containerID+":/tmp/script."+models.LanguageToExtension[language])
	return copyCmd.Run()
}

func (e *Executor) executeCodeInContainer(containerID string, language string, fileName string) (string, error) {
	execCmd := exec.Command("docker", "exec", containerID, language, fileName)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("failed to execute code in container: %w", err)
	}
	cleanupCmd := exec.Command("docker", "exec", containerID, "rm", "-f", "/tmp/script."+models.LanguageToExtension[language])
	cleanupCmd.Run()
	return string(output), nil
}
