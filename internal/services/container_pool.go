package services

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

type ContainerPool struct {
	containers chan string
	language   string
	image      string
	maxSize    int
	logger     *log.Logger
}

func (e *Executor) initializePool(pool *ContainerPool) {
	for i := 0; i < pool.maxSize; i++ {
		containerID, err := e.createContainer(pool.image)
		if err != nil {
			pool.logger.Printf("Failed to create container %d/%d for %s: %v", i+1, pool.maxSize, pool.language, err)
			continue
		}
		pool.logger.Printf("Created container %s for %s (%d/%d)", containerID[:12], pool.language, i+1, pool.maxSize)
		pool.containers <- containerID
	}
	pool.logger.Printf("Container pool for %s fully initialized with %d containers", pool.language, len(pool.containers))
}

func (e *Executor) createContainer(image string) (string, error) {
	cmd := exec.Command("docker", "run", "-d", "--net=none", "--cpus=0.5", "--memory=50m", "--entrypoint=", image, "sleep", "3600")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	containerID := strings.TrimSpace(string(output))
	return containerID, nil
}
