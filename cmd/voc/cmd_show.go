package main

import (
	"fmt"
	"strings"
	"voc/internal/i18n"
	"voc/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(showCmd)
}

var showCmd = &cobra.Command{
	Use:   "show [word]",
	Short: "Show details for a saved word",
	Run: func(cmd *cobra.Command, args []string) {
		var word string
		if len(args) > 0 {
			word = args[0]
		} else {
			// Interactive Show
			savedWords, err := db.GetAllWords()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if len(savedWords) == 0 {
				fmt.Println("No words saved.")
				return
			}

			wordList := make([]string, len(savedWords))
			for i, w := range savedWords {
				wordList[i] = w.Word
			}

			// Simple in-memory search for saved words
			searchFunc := func(query string, limit int) ([]string, error) {
				qLower := strings.ToLower(query)
				var matches []string
				for _, w := range wordList {
					if strings.Contains(strings.ToLower(w), qLower) {
						matches = append(matches, w)
					}
				}
				// Limit check
				if len(matches) > limit {
					matches = matches[:limit]
				}
				return matches, nil
			}

			previewFunc := func(w string) (string, error) {
				wd, err := db.GetWord(w)
				if err != nil {
					return "", err
				}
				if wd == nil {
					return "", nil
				}

				var lines []string
				for _, t := range wd.Types {
					lines = append(lines, t.Type+":")
					for _, d := range t.Definitions {
						lines = append(lines, "  "+d)
					}
				}
				if wd.Incomplete {
					lines = append([]string{"[NOT FOUND]"}, lines...)
				}
				return strings.Join(lines, "\n"), nil
			}

			selected, err := ui.RunFuzzyFinder(i18n.T(i18n.ShowingWord), searchFunc, previewFunc)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if selected == "" {
				fmt.Println("Cancelled.")
				return
			}
			word = selected
		}

		w, err := db.GetWord(word)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if w == nil {
			fmt.Printf("Word '%s' not found.\n", word)
			return
		}

		fmt.Println(w.Word)
		if w.Incomplete {
			fmt.Println("  [NOT FOUND]")
		}
		for _, t := range w.Types {
			fmt.Printf("  %s:\n", t.Type)
			for _, d := range t.Definitions {
				fmt.Printf("    - %s\n", d)
			}
		}

	},
}
