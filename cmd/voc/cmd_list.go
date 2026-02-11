package main

import (
	"fmt"
	"strings"
	"voc/internal/database"
	"voc/internal/i18n"
	"voc/internal/ui"

	"github.com/spf13/cobra"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved words",
	Run: func(cmd *cobra.Command, args []string) {
		words, err := db.GetAllWords()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(words) == 0 {
			fmt.Println("No words in dictionary.")
			return
		}

		var wordStrings []string
		for _, w := range words {
			wordStrings = append(wordStrings, w.Word)
		}

		// Sort words using locale-aware collation (handles accents like é with e)
		cl := collate.New(language.English, collate.Loose)
		cl.SortStrings(wordStrings)

		searchFunc := func(query string, limit int) ([]string, error) {
			var filtered []string
			q := strings.ToLower(query)
			for _, w := range wordStrings {
				if strings.Contains(strings.ToLower(w), q) {
					filtered = append(filtered, w)
				}
				if len(filtered) >= limit {
					break
				}
			}
			return filtered, nil
		}

		defFunc := func(w string) (string, error) {
			if dict == nil {
				return "", fmt.Errorf("dictionary not loaded")
			}
			return dict.Definition(w)
		}

		checkVocabFunc := func(word string) (bool, error) {
			return db.WordExists(word)
		}

		toggleVocabFunc := func(word string) (bool, error) {
			exists, err := db.WordExists(word)
			if err != nil {
				return false, err
			}

			if exists {
				if err := db.DeleteWord(word); err != nil {
					return true, err
				}
				return false, nil
			}

			// Add back to vocab
			if dict == nil {
				return false, fmt.Errorf("dictionary not loaded")
			}
			defData, err := dict.Lookup(word)
			if err != nil {
				return false, err
			}
			if defData == nil {
				return false, fmt.Errorf("%s", i18n.T(i18n.ErrWordNotFound))
			}

			var dbTypes []database.WordType
			for _, t := range defData.Types {
				dbTypes = append(dbTypes, database.WordType{
					Type:        t.Type,
					Definitions: t.Definitions,
				})
			}

			if err := db.AddWord(word, dbTypes, false); err != nil {
				return false, err
			}
			return true, nil
		}

		_, err = ui.RunFuzzyFinder(i18n.T(i18n.WordCount, len(wordStrings)), wordStrings, searchFunc, defFunc, checkVocabFunc, toggleVocabFunc)
		if err != nil {
			fmt.Printf("Err: %v\n", err)
			return
		}
	},
}
