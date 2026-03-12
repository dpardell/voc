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
	RunE: func(cmd *cobra.Command, args []string) error {
		if vocApp.Dict == nil {
			return fmt.Errorf("%s", i18n.T(i18n.DictionaryNotInstalled))
		}

		searchFunc := func(query string, limit int) ([]string, error) {
			return vocApp.Dict.Search(query, limit)
		}

		defFunc := func(w string) (string, error) {
			return vocApp.Dict.Definition(w)
		}

		checkVocabFunc := func(word string) (bool, error) {
			if vocApp.DB == nil {
				return false, nil
			}
			return vocApp.DB.WordExists(word, vocApp.TargetLang)
		}

		toggleVocabFunc := func(word string) (bool, error) {
			if vocApp.DB == nil {
				return false, nil
			}

			exists, err := vocApp.DB.WordExists(word, vocApp.TargetLang)
			if err != nil {
				return false, err
			}

			if exists {
				if err := vocApp.DB.DeleteWord(word, vocApp.TargetLang); err != nil {
					return true, err
				}
				return false, nil
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

			if err := vocApp.DB.AddWord(word, vocApp.TargetLang, dbTypes, false); err != nil {
				return false, err
			}
			return true, nil
		}

		_, err := ui.RunFuzzyFinder(i18n.T(i18n.Searching), nil, searchFunc, defFunc, checkVocabFunc, toggleVocabFunc)
		if err != nil {
			return err
		}
		return nil
	},
}
