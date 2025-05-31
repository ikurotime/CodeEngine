package services

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type ContainerPool struct {
	containers    chan string
	allContainers []string // Track all created containers
	language      string
	image         string
	maxSize       int
	logger        *log.Logger
	mu            sync.Mutex
	shutdown      bool
}

func (e *Executor) initializePool(pool *ContainerPool) {
	for i := 0; i < pool.maxSize; i++ {
		containerID, err := e.createContainer(pool.image)
		if err != nil {
			pool.logger.Printf("Failed to create container %d/%d for %s: %v", i+1, pool.maxSize, pool.language, err)
			continue
		}

		// Track the container in our list
		pool.mu.Lock()
		pool.allContainers = append(pool.allContainers, containerID)
		pool.mu.Unlock()

		pool.logger.Printf("Created container %s for %s (%d/%d)", containerID[:12], pool.language, i+1, pool.maxSize)
		pool.containers <- containerID
	}
	pool.logger.Printf("Container pool for %s fully initialized with %d containers", pool.language, len(pool.allContainers))
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

// CleanupPool stops and removes all containers in the pool
func (pool *ContainerPool) CleanupPool() {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	if pool.shutdown {
		return
	}
	pool.shutdown = true

	pool.logger.Printf("Starting cleanup for %s container pool...", pool.language)

	// Close the channel to prevent new containers from being acquired
	close(pool.containers)

	// Drain any remaining containers from the channel
	for len(pool.containers) > 0 {
		<-pool.containers
	}

	// Stop and remove ALL created containers (not just those in channel)
	pool.logger.Printf("Cleaning up %d containers for %s", len(pool.allContainers), pool.language)

	for _, containerID := range pool.allContainers {
		if err := pool.stopAndRemoveContainer(containerID); err != nil {
			pool.logger.Printf("Failed to cleanup container %s: %v", containerID[:12], err)
		} else {
			pool.logger.Printf("Successfully cleaned up container %s", containerID[:12])
		}
	}

	pool.logger.Printf("Cleanup completed for %s container pool (%d containers)", pool.language, len(pool.allContainers))
}

// stopAndRemoveContainer stops and removes a specific container
func (pool *ContainerPool) stopAndRemoveContainer(containerID string) error {
	// Stop the container (force stop after 10 seconds)
	stopCmd := exec.Command("docker", "stop", "-t", "10", containerID)
	if err := stopCmd.Run(); err != nil {
		pool.logger.Printf("Failed to stop container %s: %v", containerID[:12], err)
		// Try to force kill if stop fails
		killCmd := exec.Command("docker", "kill", containerID)
		if killErr := killCmd.Run(); killErr != nil {
			pool.logger.Printf("Failed to kill container %s: %v", containerID[:12], killErr)
		}
	}

	// Remove the container
	removeCmd := exec.Command("docker", "rm", "-f", containerID)
	if err := removeCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID[:12], err)
	}

	return nil
}

// IsShutdown returns whether the pool is in shutdown state
func (pool *ContainerPool) IsShutdown() bool {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	return pool.shutdown
}
