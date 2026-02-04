package main

import (
	"fmt"
	"voc/internal/i18n"
	"voc/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringP("comment", "c", "", "Add a comment")
}

var addCmd = &cobra.Command{
	Use:   "add [word]",
	Short: "Add a word to your dictionary",
	Run: func(cmd *cobra.Command, args []string) {
		comment, _ := cmd.Flags().GetString("comment")

		var word string
		if len(args) > 0 {
			word = args[0]
		} else {
			// Interactive mode
			if dict == nil {
				fmt.Println("Dictionary not installed. Run 'voc install-dict' first.")
				return
			}

			// Setup search and preview closures
			searchFunc := func(query string, limit int) ([]string, error) {
				return dict.Search(query, limit)
			}
			previewFunc := func(w string) (string, error) {
				return dict.Preview(w)
			}

			selected, err := ui.RunFuzzyFinder(i18n.T(i18n.AddingWord), searchFunc, previewFunc)
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

		// Add Logic
		added, err := addWordToDB(word, comment)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if !added {
			if comment == "" {
				fmt.Printf("Word '%s' already exists.\n", word)
			}
			return
		}

		printWordDetails(word)
	},
}
