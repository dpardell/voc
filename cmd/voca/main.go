package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"voca/internal/database"
	"voca/internal/dictionary"
)

var rootCmd = &cobra.Command{
	Use:   "voca",
	Short: "A personal command-line dictionary for language learners",
	Long: `Voca tracks the words you learn, fetches definitions from Wiktionary,
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
		fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
		os.Exit(1)
	}

	// Dictionary is optional for some commands, but usually needed
	dict, _ = dictionary.New()
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
