package services

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"ikurotime/code-engine/internal/models"
)

type Executor struct {
	pools    map[string]*ContainerPool
	timeout  time.Duration
	logger   *log.Logger
	mu       sync.RWMutex
	shutdown bool
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
	e.mu.RLock()
	if e.shutdown {
		e.mu.RUnlock()
		return "", fmt.Errorf("executor is shutting down")
	}
	pool, exists := e.pools[req.Language]
	e.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("unsupported language: %s", req.Language)
	}

	if pool.IsShutdown() {
		return "", fmt.Errorf("container pool for %s is shutting down", req.Language)
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
	defer os.RemoveAll(filepath.Dir(fileName))

	if err := e.copyCodeToContainer(containerID, fileName, req.Language); err != nil {
		return "", fmt.Errorf("failed to copy code to container: %w", err)
	}

	e.logger.Printf("Executing %s code in container %s", req.Language, containerID[:12])

	output, err := e.executeCodeInContainer(containerID, req.Language)

	if err != nil {
		e.logger.Printf("Code execution failed in container %s: %v", containerID[:12], err)
		return string(output), fmt.Errorf("execution failed: %w", err)
	}

	e.logger.Printf("Code execution completed successfully in container %s", containerID[:12])
	return string(output), nil
}

// Shutdown gracefully shuts down the executor and cleans up all containers
func (e *Executor) Shutdown() {
	e.mu.Lock()
	if e.shutdown {
		e.mu.Unlock()
		return
	}
	e.shutdown = true
	e.mu.Unlock()

	e.logger.Printf("Shutting down executor and cleaning up containers...")

	// Use a WaitGroup to ensure all pools are cleaned up concurrently
	var wg sync.WaitGroup
	for language, pool := range e.pools {
		wg.Add(1)
		go func(lang string, p *ContainerPool) {
			defer wg.Done()
			e.logger.Printf("Cleaning up %s container pool...", lang)
			p.CleanupPool()
		}(language, pool)
	}

	// Wait for all pools to be cleaned up
	wg.Wait()
	e.logger.Printf("All container pools cleaned up successfully")
}

// IsShutdown returns whether the executor is in shutdown state
func (e *Executor) IsShutdown() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.shutdown
}

func getContainerFromPool(pool *ContainerPool, language string) (string, error) {
	if pool.IsShutdown() {
		return "", fmt.Errorf("container pool is shutting down")
	}

	var containerID string
	select {
	case containerID = <-pool.containers:
		pool.logger.Printf("Acquired container %s for %s execution", containerID[:12], language)
		defer func() {
			// Only return container to pool if not shutting down
			if !pool.IsShutdown() {
				pool.containers <- containerID
				pool.logger.Printf("Returned container %s to %s pool", containerID[:12], language)
			}
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

func (e *Executor) executeCodeInContainer(containerID string, language string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	var execCmd *exec.Cmd
	switch language {
	case "python3":
		execCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "python3", "/tmp/script.py")
	case "nodejs":
		execCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "node", "/tmp/script.js")
	case "java":
		// Java requires compilation first
		compileCmd := exec.CommandContext(ctx, "docker", "exec", containerID, "javac", "/tmp/script.java")
		if err := compileCmd.Run(); err != nil {
			return "", fmt.Errorf("java compilation failed: %w", err)
		}
		execCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "java", "-cp", "/tmp", "script")
	case "cpp":
		// C++ requires compilation first
		compileCmd := exec.CommandContext(ctx, "docker", "exec", containerID, "g++", "/tmp/script.cpp", "-o", "/tmp/script")
		if err := compileCmd.Run(); err != nil {
			return "", fmt.Errorf("cpp compilation failed: %w", err)
		}
		execCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "/tmp/script")
	case "go":
		execCmd = exec.CommandContext(ctx, "docker", "exec", containerID, "go", "run", "/tmp/script.go")
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	output, err := execCmd.CombinedOutput()

	// Cleanup
	cleanupCmd := exec.Command("docker", "exec", containerID, "rm", "-f", "/tmp/script*")
	cleanupCmd.Run()

	if err != nil {
		return string(output), fmt.Errorf("failed to execute code in container: %w", err)
	}

	return string(output), nil
}
