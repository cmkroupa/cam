package cmd

import (
	"fmt"
	"strings"

	"cam/internal/data"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config <key> <value>",
	Short: "Set configuration values",
	Long: `Set configuration values for cam.
Supported keys:
  - model: Set the Ollama model (e.g. "qwen2.5", "llama3").
    Find models at: https://ollama.com/library`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		var value string
		if len(args) > 1 {
			value = args[1]
		}

		store := data.NewConfigStore()
		if err := store.LoadConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		switch strings.ToLower(key) {
		case "model":
			if value == "" || value == "-h" || value == "--help" {
				fmt.Printf("Current model: %s\n", store.GetOllamaModel())
				fmt.Println("Find models to use at: https://ollama.com/library")
				return nil
			}
			if err := store.SetOllamaModel(value); err != nil {
				return fmt.Errorf("failed to save model: %w", err)
			}
			fmt.Printf("Ollama model set to: %s\n", value)

		default:
			return fmt.Errorf("unknown configuration key: '%s'", key)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
