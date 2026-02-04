package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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

		for _, w := range words {
			var types []string
			for _, t := range w.Types {
				types = append(types, t.Type)
			}
			typeStr := "(no types)"
			if len(types) > 0 {
				typeStr = strings.Join(types, ", ")
			}

			marker := ""
			if w.Incomplete {
				marker = " [NOT FOUND]"
			}
			fmt.Printf("%s: %s%s\n", w.Word, typeStr, marker)
		}
	},
}
