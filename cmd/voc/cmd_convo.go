package main

import (
	"fmt"

	"voc/internal/database"
	"voc/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(convoCmd)
}

var convoCmd = &cobra.Command{
	Use:   "convo",
	Short: "Start a conversation with a language coach",
	Long:  "Practise your skills with a friendly language coach that corrects your mistakes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := vocApp.GetLLMClient()
		if err != nil {
			return fmt.Errorf("LLM error: %v (Set VERTEX_API_KEY and VERTEX_PROJECT_ID)", err)
		}
		defer client.Close()

		progress, err := database.GetProgress()
		if err != nil {
			return fmt.Errorf("error reading progress: %v", err)
		}

		if err := ui.RunConvo(client, progress); err != nil {
			return err
		}
		return nil
	},
}
