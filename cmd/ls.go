package cmd

import (
	"fmt"

	"cam/internal/data"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [Stack]",
	Short: "List all commands in a stack, or list all stacks",
	Long: `If a stack name is provided, lists all commands in that stack.
If no stack name is provided, lists all existing stacks and their item counts.

Use -p to list private stacks or commands.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		isPrivate, _ := cmd.Flags().GetBool("private")

		store := data.NewDataStore()
		if err := store.LoadData(isPrivate); err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		shouldShow := func(c data.Command) bool {
			if isPrivate {
				return c.IsPrivate
			}
			return !c.IsPrivate
		}

		if len(args) == 0 {
			if len(store.Stacks) == 0 {
				fmt.Println("No stacks found.")
				return nil
			}

			foundAny := false
			fmt.Println("Available Stacks:")
			for name, commands := range store.Stacks {
				count := 0
				for _, c := range commands {
					if shouldShow(c) {
						count++
					}
				}

				if count > 0 {
					fmt.Printf("- %s (%d commands)\n", name, count)
					foundAny = true
				}
			}

			if !foundAny {
				if isPrivate {
					fmt.Println("No private stacks found.")
				} else {
					fmt.Println("No public stacks found.")
				}
			}
			return nil
		}

		stackName := args[0]
		stack := store.GetStack(stackName)
		if len(stack) == 0 {
			if _, exists := store.Stacks[stackName]; exists {
				fmt.Printf("Stack '%s' is empty.\n", stackName)
			} else {
				fmt.Printf("Stack '%s' does not exist.\n", stackName)
			}
			return nil
		}

		fmt.Printf("Stack: %s\n", stackName)
		foundAnyCmd := false
		for i, item := range stack {
			if shouldShow(item) {
				fmt.Printf("[%d] %s\n", i, item.Cmd)
				foundAnyCmd = true
			}
		}

		if !foundAnyCmd {
			if isPrivate {
				fmt.Println("(No private commands in this stack)")
			} else {
				fmt.Println("(No public commands in this stack)")
			}
		}

		return nil
	},
}

func init() {
	lsCmd.Flags().BoolP("private", "p", false, "show private stacks/commands")
	rootCmd.AddCommand(lsCmd)
}
