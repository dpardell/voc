package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
	importCmd.Flags().StringP("file", "f", "", "Import from CSV file")
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import words from file or stdin",
	Run: func(cmd *cobra.Command, args []string) {
		filePath, _ := cmd.Flags().GetString("file")

		var input string
		if filePath != "" {
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file: %v\n", err)
				return
			}
			input = string(content)
		} else {
			// Read from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				bytes, err := io.ReadAll(os.Stdin)
				if err != nil {
					fmt.Printf("Error reading stdin: %v\n", err)
					return
				}
				input = string(bytes)
			} else {
				// No input provided
				if len(args) > 0 {
					// Try args just in case
					input = strings.Join(args, ",")
				} else {
					fmt.Println("Usage: voc import -f <file> OR pipe input")
					return
				}
			}
		}

		parts := strings.Split(input, ",")
		if strings.Contains(input, "\n") {
			parts = strings.FieldsFunc(input, func(r rune) bool {
				return r == ',' || r == '\n'
			})
		}

		count := 0
		skipped := 0
		for _, raw := range parts {
			word := strings.TrimSpace(raw)
			if word == "" {
				continue
			}

			// Use helper logic directly
			if success, _ := addWordToDB(word, ""); success {
				count++
			} else {
				skipped++
			}
		}
		fmt.Printf("Import complete: %d added, %d skipped\n", count, skipped)
	},
}
