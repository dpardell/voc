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
	db           *database.Database
	dict         *dictionary.Dictionary
	detectedLang string
)

func init() {
	// Initialize language
	lang := os.Getenv("VOC_LANG")
	if lang == "" {
		lang = os.Getenv("LANG")
	}
	if len(lang) >= 2 {
		detectedLang = lang[:2]
		i18n.SetLanguage(detectedLang)
	}

	var err error
	db, err = database.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", i18n.T(i18n.ErrDatabase), err)
		os.Exit(1)
	}

	dict, _ = dictionary.New(detectedLang)
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
