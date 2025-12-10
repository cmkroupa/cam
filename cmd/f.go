package cmd

import (
	"fmt"
	"strings"

	"cam/internal/data"

	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
)

type commandSource struct {
	items []commandItem
}

type commandItem struct {
	Stack string
	Index int
	Cmd   string
}

func (s *commandSource) String(i int) string {
	return s.items[i].Cmd
}

func (s *commandSource) Len() int {
	return len(s.items)
}

var fCmd = &cobra.Command{
	Use:   "f [query]",
	Short: "Fuzzy search to find a command across all stacks",
	Long: `Search for a command across all stacks using fuzzy matching.
Displays the stack name, index, and command for matches, sorted by relevance.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")

		store := data.NewDataStore()
		if err := store.LoadData(false); err != nil {
			return fmt.Errorf("failed to load data: %w", err)
		}

		var source commandSource
		for stackName, commands := range store.Stacks {
			for i, c := range commands {
				source.items = append(source.items, commandItem{
					Stack: stackName,
					Index: i,
					Cmd:   c.Cmd,
				})
			}
		}

		matches := fuzzy.FindFrom(query, &source)

		if len(matches) == 0 {
			fmt.Printf("No matches found for '%s'\n", query)
			return nil
		}

		for _, match := range matches {
			item := source.items[match.Index]
			fmt.Printf("[%s] [%d] %s\n", item.Stack, item.Index, item.Cmd)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(fCmd)
}
