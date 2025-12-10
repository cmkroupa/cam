package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cam",
	Short: "Camel-case Command Manager",
	Long: `cam is a CLI tool that acts as a persistent, indexed multi-clipboard
for developers to store, retrieve, and execute common cli commands.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
