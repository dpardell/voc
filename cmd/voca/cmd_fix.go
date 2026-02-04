package main

import (
	"fmt"
	"voca/internal/database"
	"voca/internal/i18n"
	"voca/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fixCmd)
}

var fixCmd = &cobra.Command{
	Use:   "fix [word]",
	Short: "Fix incomplete words",
	Run: func(cmd *cobra.Command, args []string) {
		if dict == nil {
			fmt.Println("Dictionary not installed.")
			return
		}

		var wordsToFix []string
		if len(args) > 0 {
			wordsToFix = []string{args[0]}
		} else {
			var err error
			wordsToFix, err = db.GetNotFoundWords()
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
		}

		if len(wordsToFix) == 0 {
			fmt.Println("No words to fix.")
			return
		}

		searchFunc := func(query string, limit int) ([]string, error) {
			return dict.Search(query, limit)
		}
		previewFunc := func(w string) (string, error) {
			return dict.Preview(w)
		}

		for i, w := range wordsToFix {
			fmt.Printf("[%d/%d] Fixing '%s'...\n", i+1, len(wordsToFix), w)
			selected, err := ui.RunFuzzyFinder(i18n.T(i18n.FixingWord, w), searchFunc, previewFunc)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			if selected == "" {
				fmt.Println("Skipped.")
				continue
			}

			// Apply fix
			data, err := dict.Lookup(selected)
			if err != nil || data == nil || len(data.Types) == 0 {
				fmt.Println("Error: could not find definition for selected word.")
				continue
			}

			var types []database.WordType
			for _, t := range data.Types {
				types = append(types, database.WordType{Type: t.Type, Definitions: t.Definitions})
			}

			if err := db.UpdateWordSpelling(w, selected, types); err != nil {
				fmt.Printf("Error updating word: %v\n", err)
			} else {
				fmt.Printf("Fixed: '%s' -> '%s'\n", w, selected)
			}
		}
	},
}
