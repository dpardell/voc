package main

import (
	"fmt"
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

		_, err := ui.RunFuzzyFinder(i18n.T(i18n.Searching), searchFunc, defFunc)
		if err != nil {
			fmt.Printf("Err: %v\n", err)
			return
		}
	},
}
