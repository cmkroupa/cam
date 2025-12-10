package cmd

import (
	"fmt"
	"strings"

	"cam/internal/data"

	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin <Stack> <command>",
	Short: "Pin a command to the top of a stack",
	Long: `Pin a command string to the specified stack.
New commands are prepended to the stack (index 0).

Use -p to store the command as an encrypted private command.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		commandStr := strings.Join(args[1:], " ")
		isPrivate, _ := cmd.Flags().GetBool("private")

		store := data.NewDataStore()
		if err := store.LoadData(false); err != nil {
			return fmt.Errorf("failed to load data store: %w", err)
		}

		if err := store.AddCommand(stackName, commandStr, nil, isPrivate); err != nil {
			return err
		}

		if err := store.SaveData(); err != nil {
			return fmt.Errorf("failed to save data store: %w", err)
		}

		return nil
	},
}

func init() {
	pinCmd.Flags().BoolP("private", "p", false, "encrypt command and store as private")
	rootCmd.AddCommand(pinCmd)
}
