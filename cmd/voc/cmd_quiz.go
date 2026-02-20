package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"voc/internal/llm"

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
		runAIQuiz()
	},
}

func runAIQuiz() {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: GEMINI_API_KEY environment variable not set.")
		fmt.Println("Please set it to your Gemini API key to use the AI quiz mode.")
		fmt.Println("Example: export GEMINI_API_KEY=your_key_here")
		fmt.Println("Or use --flashcards for offline mode.")
		return
	}

	client, err := llm.NewClient(apiKey)
	if err != nil {
		fmt.Printf("Error initializing AI client: %v\n", err)
		return
	}
	defer client.Close()

	fmt.Println("Generating quiz...")

	// Fetch 10 random words
	words, err := db.GetRandomWords(10)
	if err != nil {
		fmt.Printf("Error fetching words: %v\n", err)
		return
	}

	if len(words) < 4 {
		fmt.Println("Not enough words to generate a quiz. Add more words first!")
		return
	}

	var targetWords []string
	for _, w := range words {
		targetWords = append(targetWords, w.Word)
	}

	// We can use the same list for context for now
	questions, err := client.GenerateQuiz(targetWords, targetWords)
	if err != nil {
		fmt.Printf("Error generating quiz: %v\n", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	score := 0

	for i, q := range questions {
		fmt.Printf("\nQuestion %d/%d: %s\n", i+1, len(questions), q.Question)
		for j, opt := range q.Options {
			fmt.Printf("  %c) %s\n", 'A'+j, opt)
		}

		fmt.Print("Answer: ")
		if scanner.Scan() {
			ans := strings.TrimSpace(strings.ToUpper(scanner.Text()))
			if len(ans) != 1 || ans < "A" || ans > "D" {
				fmt.Println("Invalid input. Skipping.")
				continue
			}

			idx := int(ans[0] - 'A')
			if idx == q.CorrectAnswerIndex {
				fmt.Println("Correct!")
				score++
			} else {
				fmt.Printf("Wrong! The correct answer was %c) %s\n", 'A'+q.CorrectAnswerIndex, q.Options[q.CorrectAnswerIndex])
			}
		}
	}

	fmt.Printf("\nQuiz completed! Score: %d/%d\n", score, len(questions))
}

func runFlashcards() {
	words, err := db.GetAllWords()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if len(words) == 0 {
		fmt.Println("No words to quiz.")
		return
	}

	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(words))

	fmt.Println("Flashcards mode - Press Enter to see definition, Ctrl+C to exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for _, idx := range perm {
		w := words[idx]

		// Fetch full details
		fullWord, err := db.GetWord(w.Word)
		if err != nil || fullWord == nil || len(fullWord.Types) == 0 {
			continue
		}

		// Random type
		typeIdx := rand.Intn(len(fullWord.Types))
		t := fullWord.Types[typeIdx]

		fmt.Printf("Word: %s (%s)\n", fullWord.Word, t.Type)
		fmt.Print("Press Enter to see definition... ")
		scanner.Scan()

		fmt.Println("Definition:")
		for _, d := range t.Definitions {
			fmt.Printf("  - %s\n", d)
		}
		fmt.Println()
		fmt.Print("Press Enter for next word... ")
		scanner.Scan()
		fmt.Println()
	}
	fmt.Println("Quiz ended.")
}
