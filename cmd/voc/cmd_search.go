package main

import (
	"fmt"
	"voc/internal/database"
	"voc/internal/i18n"
	"voc/internal/ui"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the dictionary",
	Long:  "Search the dictionary with the option to add words to your vocabulary.",
	Run: func(cmd *cobra.Command, args []string) {
		if dict == nil {
			fmt.Println(i18n.T(i18n.DictionaryNotInstalled))
			return
		}

		searchFunc := func(query string, limit int) ([]string, error) {
			return dict.Search(query, limit)
		}
		defFunc := func(w string) (string, error) {
			return dict.Definition(w)
		}

		checkVocabFunc := func(word string) (bool, error) {
			if db == nil {
				return false, nil
			}
			return db.WordExists(word)
		}

		toggleVocabFunc := func(word string) (bool, error) {
			if db == nil {
				return false, nil
			}
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

		_, err := ui.RunFuzzyFinder(i18n.T(i18n.Searching), searchFunc, defFunc, checkVocabFunc, toggleVocabFunc)
		if err != nil {
			fmt.Printf("Err: %v\n", err)
			return
		}
	},
}
