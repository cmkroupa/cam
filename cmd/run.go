package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"cam/internal/data"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run <Stack> [index]",
	Short: "Run a command from a stack",
	Long: `Execute a command stored in a stack directly.
If no index is provided, defaults to the most recent command (index 0).
The command is executed in your default shell.`,
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
		// We need to decrypt private commands to run them
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
		if cmdStr == "" {
			return fmt.Errorf("command at index %d is empty", index)
		}

		fmt.Printf("Running: %s\n", cmdStr)

		// Determine shell to use
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}

		// Execute command interactively
		execCmd := exec.Command(shell, "-c", cmdStr)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		execCmd.Stdin = os.Stdin

		if err := execCmd.Run(); err != nil {
			// Don't duplicate the error print if the command itself failed and printed to stderr
			return fmt.Errorf("command failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
