package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"

	"voc/internal/database"
	"voc/internal/i18n"
	"voc/internal/ui"

	"github.com/spf13/cobra"
)

var flashcardsMode bool

func init() {
	quizCmd.Flags().BoolVar(&flashcardsMode, "flashcards", false, "Use legacy flashcards mode")
	rootCmd.AddCommand(quizCmd)
}

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Interactive quiz mode",
	RunE: func(cmd *cobra.Command, args []string) error {
		if flashcardsMode {
			return runFlashcards()
		}
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


func runFlashcards() error {
	words, err := vocApp.DB.GetAllWords()
	if err != nil {
		return err
	}

	if len(words) == 0 {
		return fmt.Errorf("%s", i18n.T(i18n.NoWordsToQuiz))
	}

	perm := rand.Perm(len(words))

	fmt.Println(i18n.T(i18n.FlashcardsWelcome))
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for _, idx := range perm {
		w := words[idx]

		// Fetch full details
		fullWord, err := vocApp.DB.GetWord(w.Word)
		if err != nil || fullWord == nil || len(fullWord.Types) == 0 {
			continue
		}

		// Random type
		typeIdx := rand.Intn(len(fullWord.Types))
		t := fullWord.Types[typeIdx]

		fmt.Printf("Word: %s (%s)\n", fullWord.Word, t.Type)
		fmt.Print(i18n.T(i18n.FlashcardsPress))
		scanner.Scan()

		fmt.Println("Definition:")
		for _, d := range t.Definitions {
			fmt.Printf("  - %s\n", d)
		}
		fmt.Println()
		fmt.Print(i18n.T(i18n.FlashcardsNext))
		scanner.Scan()
		fmt.Println()
	}
	fmt.Println(i18n.T(i18n.FlashcardsEnd))
	return nil
}
