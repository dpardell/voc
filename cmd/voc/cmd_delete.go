package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().Bool("all", false, "Delete all words")
}

var deleteCmd = &cobra.Command{
	Use:   "delete [word]",
	Short: "Delete a word",
	Run: func(cmd *cobra.Command, args []string) {
		deleteAll, _ := cmd.Flags().GetBool("all")
		if deleteAll {
			fmt.Printf("WARNING: This will delete ALL words. Type 'yes' to confirm: ")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			if strings.TrimSpace(strings.ToLower(scanner.Text())) == "yes" {
				if err := db.DeleteAllWords(); err != nil {
					fmt.Printf("Error: %v\n", err)
				} else {
					fmt.Println("All words deleted.")
				}
			} else {
				fmt.Println("Cancelled.")
			}
			return
		}

		if len(args) == 0 {
			fmt.Println("Error: word required")
			return
		}

		word := args[0]
		if err := db.DeleteWord(word); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Deleted word: %s\n", word)
		}
	},
}
