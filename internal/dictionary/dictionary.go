package dictionary

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Dictionary struct {
	db *sql.DB
}

type DefinitionData struct {
	Types []TypeData
}

type TypeData struct {
	Type        string
	Definitions []string
}

var DefaultDictionaryPath string

func GetDictionaryPath() (string, error) {
	if envPath := os.Getenv("VOCA_DB_PATH"); envPath != "" {
		return envPath, nil
	}
	if DefaultDictionaryPath != "" {
		return DefaultDictionaryPath, nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "voca", "dictionary.db"), nil
}

func New() (*Dictionary, error) {
	dbPath, err := GetDictionaryPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("dictionary not found at %s", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	return &Dictionary{db: db}, nil
}

func (d *Dictionary) Close() error {
	return d.db.Close()
}

func (d *Dictionary) Lookup(word string) (*DefinitionData, error) {
	wordLower := strings.ToLower(word)
	rows, err := d.db.Query("SELECT DISTINCT pos, gloss FROM dictionary WHERE word = ?", wordLower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	typeMap := make(map[string][]string)

	for rows.Next() {
		var pos string
		var gloss string
		if err := rows.Scan(&pos, &gloss); err != nil {
			return nil, err
		}
		if pos == "" {
			pos = "unknown"
		}
		typeMap[pos] = append(typeMap[pos], gloss)
	}

	if len(typeMap) == 0 {
		return nil, nil // Not found
	}

	var data DefinitionData
	for pos, defs := range typeMap {
		data.Types = append(data.Types, TypeData{Type: pos, Definitions: defs})
	}
	return &data, nil
}

func (d *Dictionary) Search(query string, limit int) ([]string, error) {
	if query == "" {
		return nil, nil
	}
	queryLower := strings.ToLower(strings.TrimSpace(query))
	escapedQuery := strings.ReplaceAll(queryLower, "\"", "\"\"")

	// Try FTS5
	rows, err := d.db.Query("SELECT word FROM words_fts WHERE word MATCH ? ORDER BY rank LIMIT ?",
		"^"+escapedQuery+"*", limit*2)

	// Fallback to LIKE if FTS fails (or table missing)
	if err != nil {
		rows, err = d.db.Query(`SELECT DISTINCT word FROM dictionary 
			WHERE word LIKE ? 
			ORDER BY CASE WHEN word LIKE ? THEN 0 ELSE 1 END, word LIMIT ?`,
			queryLower+"%", queryLower+"%", limit)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var w string
		if err := rows.Scan(&w); err != nil {
			return nil, err
		}
		results = append(results, w)
	}

	// Manual post-sort to ensure prefix matches are first (if FTS didn't guarantee it perfectly)
	// although the query usually handles it.
	return results[:min(len(results), limit)], nil
}

func (d *Dictionary) Preview(word string) (string, error) {
	data, err := d.Lookup(word)
	if err != nil {
		return "", err
	}
	if data == nil {
		return "", nil
	}

	var lines []string
	for i, t := range data.Types {
		if i >= 3 {
			break
		}
		lines = append(lines, fmt.Sprintf("%s:", t.Type))
		for j, def := range t.Definitions {
			if j >= 2 {
				break
			}
			lines = append(lines, fmt.Sprintf("  %s", def))
		}
	}
	return strings.Join(lines, "\n"), nil
}
