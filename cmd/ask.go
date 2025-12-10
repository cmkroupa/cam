package cmd

import (
	"context"
	"fmt"
	"strings"

	"cam/internal/data"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask Gemini anything",
	Long: `Ask Gemini anything.
Requires a Gemini API key to be set via 'cam config api-key <KEY>'.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		question := strings.Join(args, " ")

		configStore := data.NewConfigStore()
		if err := configStore.LoadConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		apiKey := configStore.GetAPIKey()
		if apiKey == "" {
			return fmt.Errorf("API key not set. Please run 'cam config api-key <YOUR_KEY>'")
		}

		ctx := context.Background()
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			return fmt.Errorf("failed to create Gemini client: %w", err)
		}

		model := "gemini-2.5-flash"
		var prompt string
		oneline, _ := cmd.Flags().GetBool("oneline")
		if oneline {
			prompt = fmt.Sprintf("You are a command line expert. If the user asks how to do something in the terminal, return ONLY the command string. If the user asks a general question, return the answer as a raw string. Do NOT wrap the answer in 'echo' unless explicitly asked. Do not use markdown. Question: %s", question)
		} else {
			prompt = fmt.Sprintf("You are a command line expert. Provide a VERY CONCISE, readable answer to the following question, using markdown code blocks where appropriate: %s", question)
		}

		resp, err := client.Models.GenerateContent(ctx, model, genai.Text(prompt), nil)
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}

		if resp == nil {
			return fmt.Errorf("no response from Gemini")
		}

		var resultText string
		if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				if part.Text != "" {
					resultText += part.Text
				}
			}
		}

		if resultText == "" {
			return fmt.Errorf("empty response from Gemini")
		}

		if !oneline {
			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
				glamour.WithWordWrap(80),
			)
			if err != nil {
				return fmt.Errorf("failed to create markdown renderer: %w", err)
			}

			out, err := renderer.Render(resultText)
			if err != nil {
				return fmt.Errorf("failed to render markdown: %w", err)
			}
			fmt.Print(out)
		} else {
			fmt.Println(strings.TrimSpace(resultText))
		}
		return nil
	},
}

func init() {
	askCmd.Flags().BoolP("oneline", "o", false, "Print only the command string (no markdown)")
	rootCmd.AddCommand(askCmd)
}
