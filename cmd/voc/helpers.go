package main

import (
	"fmt"
	"voc/internal/database"
	"voc/internal/i18n"
)

func addWordToDB(word string) (bool, error) {
	exists, err := db.WordExists(word)
	if err != nil {
		return false, fmt.Errorf("database error: %v", err)
	}

	if exists {
		return false, nil
	}

	wordTypes := []database.WordType{}
	incomplete := true

	if dict != nil {
		data, err := dict.Lookup(word)
		if err != nil {
			return false, fmt.Errorf("Err: %v", err)
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

	if incomplete {
		return false, fmt.Errorf("%s", i18n.T(i18n.ErrWordNotFound))
	}

	if err := db.AddWord(word, wordTypes, incomplete); err != nil {
		return false, fmt.Errorf("Err: %v", err)
	}

	return true, nil
}
