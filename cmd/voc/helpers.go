package main

import (
	"fmt"
	"voc/internal/database"
	"voc/internal/dictionary"
)

// addWordToDB handles the complete flow of adding a word:
// 1. Checks if it exists (returns if so)
// 2. Looks up definition in dictionary
// 3. Adds to DB
// 4. Adds comment if provided
// Returns true if word was added, false if it already existed
func addWordToDB(word string, comment string) (bool, error) {
	exists, err := db.WordExists(word)
	if err != nil {
		return false, fmt.Errorf("database error: %v", err)
	}

	if exists {
		if comment != "" {
			if err := db.AddComment(word, comment); err != nil {
				return false, fmt.Errorf("error adding comment: %v", err)
			}
			fmt.Printf("Word '%s' already exists. Comment added.\n", word)
		} else {
			// For import command, we might not want to print this every time?
			// The original logic printed "already exists" for 'add', but silently skipped or counted skipped in 'import'.
			// I'll make this helper focused on the ACTION. The caller can handle specific UI feedback if needed,
			// but for simple re-use, I'll print if it's a single add.
			// Actually, for 'import', we want to return 'false' so it increments 'skipped'.
		}
		return false, nil
	}

	wordTypes := []database.WordType{}
	incomplete := true

	if dict != nil {
		data, err := dict.Lookup(word)
		if err != nil {
			// Log error but continue?
			fmt.Printf("Dictionary look up error: %v\n", err)
		}

		if data != nil && len(data.Types) > 0 {
			incomplete = false
			for _, t := range data.Types {
				wt := database.WordType{
					Type:        t.Type,
					Definitions: t.Definitions,
				}
				wordTypes = append(wordTypes, wt)
			}
		}
	}

	if err := db.AddWord(word, wordTypes, incomplete); err != nil {
		return false, fmt.Errorf("error adding word: %v", err)
	}

	if comment != "" {
		if err := db.AddComment(word, comment); err != nil {
			fmt.Printf("Error adding comment: %v\n", err)
		}
	}

	// Return the added word object or details?
	// For "add" command we want to print details. For "import" we just want to know it succeeded.

	return true, nil
}

func printWordDetails(word string) {
	w, err := db.GetWord(word)
	if err != nil || w == nil {
		return
	}

	if w.Incomplete {
		fmt.Printf("Added: %s (no definition found)\n", w.Word)
	} else {
		fmt.Printf("Added: %s\n", w.Word)
		for _, t := range w.Types {
			fmt.Printf("  %s:\n", t.Type)
			for _, d := range t.Definitions {
				fmt.Printf("    - %s\n", d)
			}
		}
	}
}

func printDefinition(word string, data *dictionary.DefinitionData) {
	if data == nil {
		fmt.Printf("No definition found for '%s'\n", word)
		return
	}

	fmt.Printf("Word: %s\n", word)
	for _, t := range data.Types {
		fmt.Printf("  %s:\n", t.Type)
		for _, d := range t.Definitions {
			fmt.Printf("    - %s\n", d)
		}
	}
}
