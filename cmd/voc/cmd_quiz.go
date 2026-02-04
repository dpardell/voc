package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(quizCmd)
}

var quizCmd = &cobra.Command{
	Use:   "quiz",
	Short: "Interactive quiz mode",
	Run: func(cmd *cobra.Command, args []string) {
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

		fmt.Println("Quiz mode - Press Enter to see definition, Ctrl+C to exit")
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
	},
}
