package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.Flags().StringP("add", "a", "", "Add a comment")
	commentCmd.Flags().IntP("delete", "d", 0, "Delete comment by ID")
}

var commentCmd = &cobra.Command{
	Use:   "comment [word]",
	Short: "Manage comments",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: word required")
			return
		}
		word := args[0]

		addMsg, _ := cmd.Flags().GetString("add")
		delID, _ := cmd.Flags().GetInt("delete")

		if addMsg != "" {
			if err := db.AddComment(word, addMsg); err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println("Comment added.")
			}
		} else if delID != 0 {
			wordID, err := db.GetWordID(word)
			if err != nil {
				fmt.Printf("Error: Word '%s' not found.\n", word)
				return
			}

			if err := db.DeleteComment(wordID, delID); err != nil {
				fmt.Printf("Error deleting comment: %v\n", err)
			} else {
				fmt.Printf("Comment %d deleted.\n", delID)
			}
		} else {
			// Show
			comments, err := db.GetComments(word)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			if len(comments) == 0 {
				fmt.Printf("No comments for '%s'\n", word)
				return
			}
			fmt.Printf("Comments for '%s':\n", word)
			for _, c := range comments {
				fmt.Printf("  %d. %s\n", c.ID, c.Comment)
			}
		}
	},
}
