package main

import (
	"fmt"
	"voca/internal/i18n"
	"voca/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search [word]",
	Short: "Search for a word in the dictionary",
	Long:  "Search for a word in the dictionary without adding it to your vocabulary list. Supports interactive search if no word is provided.",
	Run: func(cmd *cobra.Command, args []string) {
		if dict == nil {
			fmt.Println("Dictionary not installed. Run 'wordsoup install-dict' first.")
			return
		}

		var word string
		if len(args) > 0 {
			word = args[0]
		} else {
			// Interactive mode
			searchFunc := func(query string, limit int) ([]string, error) {
				return dict.Search(query, limit)
			}
			previewFunc := func(w string) (string, error) {
				return dict.Preview(w)
			}

			selected, err := ui.RunFuzzyFinder(i18n.T(i18n.Searching), searchFunc, previewFunc)
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

		// Lookup and print definition
		data, err := dict.Lookup(word)
		if err != nil {
			fmt.Printf("Error looking up word: %v\n", err)
			return
		}

		printDefinition(word, data)
	},
}
