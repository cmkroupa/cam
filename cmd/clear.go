package cmd

import (
	"cam/internal/data"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear <Stack> [index]",
	Short: "Remove a stack or a specific command from a stack",
	Long: `Remove a stack or a specific command.
To remove a specific command, provide the stack name and the index.
To remove an entire stack, provide only the stack name.
Use -a to remove ALL stacks.`,
	Args: func(cmd *cobra.Command, args []string) error {
		deleteAll, _ := cmd.Flags().GetBool("all")
		if deleteAll {
			return cobra.NoArgs(cmd, args)
		}
		return cobra.RangeArgs(1, 2)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		deleteAll, _ := cmd.Flags().GetBool("all")

		store := data.NewDataStore()
		if err := store.LoadData(true); err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		if deleteAll {
			// Clear everything
			store.Stacks = make(map[string][]data.Command)
			if err := store.SaveData(); err != nil {
				return fmt.Errorf("failed to save data: %w", err)
			}
			return nil
		}

		stackName := args[0]

		if len(args) == 2 {
			index, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid index: %s", args[1])
			}

			if err := store.RemoveCommand(stackName, index); err != nil {
				return err
			}
		} else {
			if err := store.RemoveStack(stackName); err != nil {
				return err
			}
		}

		if err := store.SaveData(); err != nil {
			return fmt.Errorf("failed to save data: %w", err)
		}

		return nil
	},
}

func init() {
	clearCmd.Flags().BoolP("all", "a", false, "remove all stacks")
	rootCmd.AddCommand(clearCmd)
}
