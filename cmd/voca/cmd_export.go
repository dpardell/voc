package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP("file", "f", "", "Export to file")
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export words to file or stdout",
	Run: func(cmd *cobra.Command, args []string) {
		words, err := db.GetAllWords()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		var wordList []string
		for _, w := range words {
			wordList = append(wordList, w.Word)
		}
		output := strings.Join(wordList, ", ")

		filePath, _ := cmd.Flags().GetString("file")
		if filePath != "" {
			err := os.WriteFile(filePath, []byte(output), 0644)
			if err != nil {
				fmt.Printf("Error writing file: %v\n", err)
			} else {
				fmt.Printf("Exported %d words to %s\n", len(words), filePath)
			}
		} else {
			fmt.Println(output)
		}
	},
}
