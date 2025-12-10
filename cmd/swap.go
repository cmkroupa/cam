package cmd

import (
	"fmt"
	"strconv"

	"cam/internal/data"

	"github.com/spf13/cobra"
)

var swapCmd = &cobra.Command{
	Use:   "swap <Stack> <index> [index = 0]",
	Short: "Swap the position of two commands",
	Long: `Swap the position of two commands in a stack. 
	With a default swap position of to the first index.

Example:
  cam swap python 2
  cam swap python 2 1`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		stackName := args[0]
		index1, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid index: %s", args[1])
		}

		index2 := 0
		if len(args) == 3 {
			index2, err = strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid index: %s", args[2])
			}
		}

		store := data.NewDataStore()
		if err := store.LoadData(false); err != nil {
			return fmt.Errorf("failed to load data store: %w", err)
		}

		if err := store.Swap(stackName, index1, index2); err != nil {
			return err
		}

		if err := store.SaveData(); err != nil {
			return fmt.Errorf("failed to save data store: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(swapCmd)
}
