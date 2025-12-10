package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cam/internal/ai"
	"cam/internal/data"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask Gemini anything",
	Long: `Ask your local Ollama model anything.
Requires Ollama to be installed and running.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := strings.Join(args, " ")

		configStore := data.NewConfigStore()
		if err := configStore.LoadConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		ctx := context.Background()

		oneline, _ := cmd.Flags().GetBool("oneline")
		withContext, _ := cmd.Flags().GetBool("context")

		var contextStr string
		if withContext {
			contextStr = "directory tree:\n.\n" + buildFileTree(".", "", 0, 2)
		}

		var prompt strings.Builder
		prompt.WriteString("SYSTEM: You are a concise CLI technical assistant.\n")
		prompt.WriteString("RULES:\n")
		prompt.WriteString("1. Keep answers short, accurate, and direct.\n")
		prompt.WriteString("2. Avoid conversational filler (e.g. 'Here is a summary', 'I hope this helps').\n")
		prompt.WriteString("3. Use markdown code blocks for examples.\n")
		prompt.WriteString("\n")

		if contextStr != "" {
			prompt.WriteString(fmt.Sprintf("CONTEXT (Current Directory):\n%s\n", contextStr))
		}

		if oneline {
			prompt.WriteString("Provide a single-line plain text answer. No markdown. No explanations.\n")
		} else {
			prompt.WriteString("Provide a concise explanation using markdown.\n")
		}

		prompt.WriteString(fmt.Sprintf("USER QUESTION: %s", question))

		resultText, err := ai.GenerateContent(ctx, configStore, prompt.String())
		if err != nil {
			return err
		}

		if oneline {
			fmt.Println(strings.TrimSpace(resultText))
		} else {
			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(80),
			)
			if err != nil {
				return fmt.Errorf("failed to create markdown renderer: %w", err)
			}

			out, err := renderer.Render(resultText)
			if err != nil {
				return fmt.Errorf("failed to render markdown: %w", err)
			}
			fmt.Print(out)
		}
		return nil
	},
}

func init() {
	askCmd.Flags().BoolP("oneline", "o", false, "Get a concise one-line answer")
	askCmd.Flags().BoolP("context", "c", false, "Include local file context")
	rootCmd.AddCommand(askCmd)
}

// Shared with cmdr.go (package scope)

func buildFileTree(dir string, prefix string, depth int, maxDepth int) string {
	if depth > maxDepth {
		return ""
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	var validEntries []os.DirEntry
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		if e.IsDir() && ignoredDirs[name] {
			continue
		}
		validEntries = append(validEntries, e)
	}

	var sb strings.Builder
	for i, e := range validEntries {
		isLast := i == len(validEntries)-1
		connector := "├── "
		newPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		name := e.Name()
		if e.IsDir() {
			name += "/"
		}

		sb.WriteString(fmt.Sprintf("%s%s%s\n", prefix, connector, name))

		if e.IsDir() {
			sb.WriteString(buildFileTree(filepath.Join(dir, e.Name()), newPrefix, depth+1, maxDepth))
		}
	}
	return sb.String()
}
