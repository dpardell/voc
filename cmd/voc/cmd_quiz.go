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
	Run: func(cmd *cobra.Command, args []string) {
		if flashcardsMode {
			runFlashcards()
			return
		}
		runAIQuiz(cmd.Context())
	},
}

func runAIQuiz(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	client, err := vocApp.GetLLMClient()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Println("Please set VERTEX_API_KEY and VERTEX_PROJECT_ID.")
		fmt.Println("Or use --flashcards for offline mode.")
		fmt.Println("\n" + i18n.T(i18n.PressEnterToCont))
		bufio.NewScanner(os.Stdin).Scan()
		return
	}
	defer client.Close()

	progress, err := database.GetProgress()
	if err != nil {
		fmt.Printf("Error reading progress: %v\n", err)
		return
	}

	// Fetch 10 random words for the quiz
	words, err := vocApp.DB.GetRandomWords(10)
	if err != nil {
		fmt.Printf("Error fetching words: %v\n", err)
		return
	}

	if len(words) < 4 {
		fmt.Println(i18n.T(i18n.ErrorNoWords))
		fmt.Println("\n" + i18n.T(i18n.PressEnterToCont))
		bufio.NewScanner(os.Stdin).Scan()
		return
	}

	var targetWords []string
	for _, w := range words {
		targetWords = append(targetWords, w.Word)
	}

	result, err := ui.RunQuiz(client, targetWords, progress)
	if err != nil {
		fmt.Printf("Error running quiz: %v\n", err)
		return
	}

	if result != nil && result.NewProgress != "" {
		progressPath, _ := database.GetProgressPath()
		if err := os.WriteFile(progressPath, []byte(result.NewProgress), 0644); err != nil {
			fmt.Printf("Error saving progress file: %v\n", err)
		}
	}
}


func runFlashcards() {
	words, err := vocApp.DB.GetAllWords()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(words) == 0 {
		fmt.Println(i18n.T(i18n.NoWordsToQuiz))
		return
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
}
