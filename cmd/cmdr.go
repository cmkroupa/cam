package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cam/internal/ai"
	"cam/internal/data"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	commandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")). // Pink/Magentaish
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")). // Purpleish
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Cyan/Greenish
			Bold(true)
)

var cmdrCmd = &cobra.Command{
	Use:   "cmdr [question]",
	Short: "Generate a shell command from a question",
	Long: `Ask your local Ollama model to generate a specific shell command based on your request.
Returns ONLY the command string, ready to copy-paste or pipe.
Requires Ollama to be installed and running.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := strings.Join(args, " ")

		configStore := data.NewConfigStore()
		if err := configStore.LoadConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		ctx := context.Background() // Context is always included for cmdr

		contextStr := "Directory structure (flat list):\n" + getFlatFileList(".", 3)
		copyToClipboard, _ := cmd.Flags().GetBool("copy")

		var prompt strings.Builder

		prompt.WriteString("SYSTEM: You are a command-line interface expert. Your goal is to provide the exact shell command(s) the user needs.\n")
		prompt.WriteString("OBJECTIVE: Convert the user's request (which might be a question or a statement) into a single valid shell command line.\n")
		prompt.WriteString("RULES:\n")
		prompt.WriteString("1. Output ONLY the command text. Do not include markdown formatting (like ```bash). Do not include explanations.\n")
		prompt.WriteString("2. If multiple steps are required, chain them using '&&' or ';'.\n")
		prompt.WriteString("3. If the user asks 'how to' or 'steps to', provide the actual commands to perform those steps.\n")
		prompt.WriteString("4. Use the provided file list to resolve paths if applicable.\n")
		prompt.WriteString("5. Assume a modern shell (bash/zsh).\n")
		prompt.WriteString("\n")

		if contextStr != "" {
			prompt.WriteString(fmt.Sprintf("CONTEXT (File List):\n%s\n", contextStr))
			prompt.WriteString("CRITICAL: Use the paths above to correct the user's request if needed.\n")
		}

		prompt.WriteString(fmt.Sprintf("USER REQUEST: %s\n", question))
		prompt.WriteString("COMMAND:") // Pre-fill the start to encourage completion

		resultText, err := ai.GenerateContent(ctx, configStore, prompt.String())
		if err != nil {
			return err
		}

		// Intelligent Parsing
		// 1. If markdown blocks exist, take the content inside the first one
		if start := strings.Index(resultText, "```"); start != -1 {
			// Find end
			end := strings.Index(resultText[start+3:], "```")
			if end != -1 {
				// content is between start+3 and start+3+end
				codeBlock := resultText[start+3 : start+3+end]
				// remove optional language identifier like "bash" or "sh"
				lines := strings.Split(codeBlock, "\n")
				if len(lines) > 0 {
					firstLine := strings.TrimSpace(lines[0])
					if firstLine == "bash" || firstLine == "sh" || firstLine == "zsh" {
						codeBlock = strings.Join(lines[1:], "\n")
					}
				}
				resultText = codeBlock
			}
		}

		// Intelligent Cleanup (Pre-Processing)
		resultText = strings.TrimSpace(resultText)

		// 2. Intelligent Cleanup: Remove standalone language identifiers or comments if they appear at the start
		lines := strings.Split(resultText, "\n")

		// Helper to check if a line is a language tag
		isLanguageTag := func(s string) bool {
			s = strings.TrimSpace(strings.ToLower(s))
			switch s {
			case "bash", "sh", "zsh", "shell", "console":
				return true
			}
			return false
		}

		// Strip leading lines that are just language tags or comments
		for len(lines) > 0 {
			first := strings.TrimSpace(lines[0])
			if isLanguageTag(first) || strings.HasPrefix(first, "#") {
				lines = lines[1:] // drop the line
				continue
			}
			break
		}

		if len(lines) == 0 {
			return fmt.Errorf("AI response contained only filtered lines (comments/tags). Raw: %s", resultText)
		}

		resultText = strings.Join(lines, "\n")
		resultText = strings.TrimSpace(resultText)

		// 2. If no markdown, but multiple lines, try to find the one that looks like a command?
		// For now, prompt engineering should fix this, but let's just trim.

		if resultText == "" {
			return fmt.Errorf("empty response from AI")
		}

		if copyToClipboard {
			if err := clipboard.WriteAll(resultText); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to copy to clipboard: %v\n", err)
			} else {
				fmt.Println(successStyle.Render("✔ Copied to clipboard!"))
			}
		}

		fmt.Println(commandStyle.Render(resultText))
		return nil
	},
}

var ignoredDirs = map[string]bool{
	"node_modules": true,
	"vendor":       true,
	"dist":         true,
	"build":        true,
	".git":         true,
	".tea":         true,
	"__pycache__":  true,
}

func getFlatFileList(root string, maxDepth int) string {
	var sb strings.Builder

	// Helper for recursion
	var walk func(path string, depth int)
	walk = func(path string, depth int) {
		if depth > maxDepth {
			return
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return
		}

		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, ".") {
				continue // skip hidden
			}
			if e.IsDir() && ignoredDirs[name] {
				continue
			}

			fullPath := filepath.Join(path, name)
			sb.WriteString(fullPath + "\n")

			if e.IsDir() {
				walk(fullPath, depth+1)
			}
		}
	}

	walk(root, 1)
	return sb.String()
}

func init() {
	cmdrCmd.Flags().BoolP("copy", "c", false, "Copy generated command to clipboard")
	rootCmd.AddCommand(cmdrCmd)
}
