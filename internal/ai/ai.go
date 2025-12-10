package ai

import (
	"cam/internal/data"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// GenerateContent abstracts the AI provider to return text from a prompt
func GenerateContent(ctx context.Context, configStore *data.ConfigStore, prompt string) (string, error) {
	model := configStore.GetOllamaModel()
	return generateOllamaRun(ctx, model, prompt)
}

func generateOllamaRun(ctx context.Context, model string, prompt string) (string, error) {
	// We use 'ollama run <model> <prompt>' directly.
	// This handles auto-pulling and execution.
	cmd := exec.CommandContext(ctx, "ollama", "run", model, prompt)

	// Stream stderr to user so they see "pulling..." or errors, but don't capture it in output
	cmd.Stderr = os.Stderr

	output, err := cmd.Output()
	if err != nil {
		// If error, the user likely saw it on stderr, but we return a generic error or the error from exec
		return "", fmt.Errorf("ollama run failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}
