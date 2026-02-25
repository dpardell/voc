package main

import (
	"fmt"
	"strings"
	"voc/internal/database"
	"voc/internal/i18n"
	"voc/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved words",
	RunE: func(cmd *cobra.Command, args []string) error {
		words, err := vocApp.DB.GetAllWords()
		if err != nil {
			return err
		}

		if len(words) == 0 {
			return fmt.Errorf("%s", i18n.T(i18n.ErrorNoSavedWords))
		}

		var wordStrings []string
		for _, w := range words {
			wordStrings = append(wordStrings, w.Word)
		}

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
			if vocApp.Dict == nil {
				return "", fmt.Errorf("%s", i18n.T(i18n.DictionaryLoading))
			}
			return vocApp.Dict.Definition(w)
		}

		checkVocabFunc := func(word string) (bool, error) {
			return vocApp.DB.WordExists(word)
		}

		toggleVocabFunc := func(word string) (bool, error) {
			exists, err := vocApp.DB.WordExists(word)
			if err != nil {
				return false, err
			}

			if exists {
				if err := vocApp.DB.DeleteWord(word); err != nil {
					return true, err
				}
				return false, nil
			}

			// Add back to vocab
			if vocApp.Dict == nil {
				return false, fmt.Errorf("%s", i18n.T(i18n.DictionaryLoading))
			}
			defData, err := vocApp.Dict.Lookup(word)
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

			if err := vocApp.DB.AddWord(word, dbTypes, false); err != nil {
				return false, err
			}
			return true, nil
		}

		_, err = ui.RunFuzzyFinder(i18n.T(i18n.WordCount, len(wordStrings)), wordStrings, searchFunc, defFunc, checkVocabFunc, toggleVocabFunc)
		if err != nil {
			return err
		}
		return nil
	},
}
