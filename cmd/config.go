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
Currently supported keys:
  - api-key: Set the Gemini API key.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		store := data.NewConfigStore()
		if err := store.LoadConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		switch strings.ToLower(key) {
		case "api-key":
			if err := store.SetAPIKey(value); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

		default:
			return fmt.Errorf("unknown configuration key: '%s'", key)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
