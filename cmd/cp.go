package cmd

import (
	"fmt"
	"strconv"

	"cam/internal/data"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var cpCmd = &cobra.Command{
	Use:   "cp <Stack> [index]",
	Short: "Copy a command from the stack to the system clipboard",
	Long: `Copy a command to the clipboard.
If no index is provided, defaults to the most recent command (index 0).`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		index := 0

		if len(args) == 2 {
			var err error
			index, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid index provided: %s (must be an integer)", args[1])
			}
		}

		store := data.NewDataStore()
		if err := store.LoadData(true); err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		stack := store.GetStack(stackName)
		if len(stack) == 0 {
			return fmt.Errorf("stack '%s' is empty or does not exist", stackName)
		}

		if index < 0 || index >= len(stack) {
			return fmt.Errorf("index %d is out of bounds for stack '%s' (length %d)", index, stackName, len(stack))
		}

		cmdStr := stack[index].Cmd

		if err := clipboard.WriteAll(cmdStr); err != nil {
			return fmt.Errorf("failed to write to clipboard: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
}
