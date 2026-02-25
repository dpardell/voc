package main

import (
	"fmt"

	"voc/internal/database"
	"voc/internal/i18n"
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
	Run: func(cmd *cobra.Command, args []string) {
		client, err := vocApp.GetLLMClient()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Println("\nPress Enter to continue...")
			var discard string
			fmt.Scanln(&discard)
			return
		}
		defer client.Close()

		progress, err := database.GetProgress()
		if err != nil {
			fmt.Printf("Error reading progress: %v\n", err)
			return
		}

		if err := ui.RunConvo(client, progress); err != nil {
			fmt.Printf("Error running conversation: %v\n", err)
			fmt.Println("\n" + i18n.T(i18n.PressEnterToCont))
			var discard string
			fmt.Scanln(&discard)
		}
	},
}
