package main

import (
	"context"
	"fmt"
	"os"

	"voc/internal/database"
	"voc/internal/i18n"
	"voc/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(quizCmd)
}

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Interactive quiz mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAIQuiz(cmd.Context())
	},
}

func runAIQuiz(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	client, err := vocApp.GetLLMClient()
	if err != nil {
		return fmt.Errorf("LLM error: %v (Set VERTEX_API_KEY and VERTEX_PROJECT_ID)", err)
	}
	defer client.Close()

	progress, err := database.GetProgress()
	if err != nil {
		return fmt.Errorf("error reading progress: %v", err)
	}

	// Fetch 10 random words for the quiz
	words, err := vocApp.DB.GetRandomWords(10)
	if err != nil {
		return fmt.Errorf("error fetching words: %v", err)
	}

	if len(words) < 4 {
		return fmt.Errorf("%s", i18n.T(i18n.ErrorNoWords))
	}

	var targetWords []string
	for _, w := range words {
		targetWords = append(targetWords, w.Word)
	}

	result, err := ui.RunQuiz(client, targetWords, progress)
	if err != nil {
		return fmt.Errorf("error running quiz: %v", err)
	}

	if result != nil && result.NewProgress != "" {
		progressPath, _ := database.GetProgressPath()
		if err := os.WriteFile(progressPath, []byte(result.NewProgress), 0644); err != nil {
			return fmt.Errorf("error saving progress file: %v", err)
		}
	}
	return nil
}
