package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"voc/internal/database"
	"voc/internal/dictionary"
	"voc/internal/i18n"
)

var rootCmd = &cobra.Command{
	Use:   "voc",
	Short: "Voc - Vocabulary trainer",
	Long: `Voc tracks the words you learn, fetches definitions from Wiktionary,
and helps you review with interactive quizzes.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Global instances
var (
	db   *database.Database
	dict *dictionary.Dictionary
)

func init() {
	var err error
	db, err = database.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", i18n.T(i18n.ErrDatabase), err)
		os.Exit(1)
	}

	dict, _ = dictionary.New(i18n.GetLanguage())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
