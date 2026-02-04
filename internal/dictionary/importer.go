package dictionary

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	KaikkiURL = "https://kaikki.org/frwiktionary/Fran%C3%A7ais/kaikki.org-dictionary-Fran%C3%A7ais.jsonl.gz"
)

type Importer struct {
	DictDir   string
	DictDB    string
	CacheFile string
}

func NewImporter() (*Importer, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dictDir := filepath.Join(home, ".local", "share", "voca")
	if err := os.MkdirAll(dictDir, 0755); err != nil {
		return nil, err
	}

	return &Importer{
		DictDir:   dictDir,
		DictDB:    filepath.Join(dictDir, "dictionary.db"),
		CacheFile: filepath.Join(dictDir, "kaikki-francais.jsonl"),
	}, nil
}

func (i *Importer) DownloadAndImport(force bool) error {
	if _, err := os.Stat(i.CacheFile); err == nil && !force {
		fmt.Println("Using cached dictionary file...")
		fmt.Println("  (use --force to re-download)")
		return i.importFile()
	}

	fmt.Println("Downloading French dictionary from kaikki.org...")

	resp, err := http.Get(KaikkiURL)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Direct piping from response body to decompressor to file
	fmt.Println("Decompressing and saving...")
	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	outFile, err := os.Create(i.CacheFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzReader)
	if err != nil {
		return err
	}

	return i.importFile()
}

func (i *Importer) importFile() error {
	fmt.Println("Importing into database...")

	os.Remove(i.DictDB) // Clean slate

	db, err := sql.Open("sqlite3", i.DictDB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Initial schema
	if _, err := db.Exec(`
		CREATE TABLE dictionary (word TEXT NOT NULL, pos TEXT, gloss TEXT NOT NULL);
		CREATE INDEX idx_word ON dictionary(word);
		CREATE VIRTUAL TABLE words_fts USING fts5(word, tokenize='unicode61');
	`); err != nil {
		return err
	}

	file, err := os.Open(i.CacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Larger buffer for long lines
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO dictionary (word, pos, gloss) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	ftsStmt, err := tx.Prepare("INSERT INTO words_fts (word) VALUES (?)")
	if err != nil {
		return err
	}
	defer ftsStmt.Close()

	uniqueWords := make(map[string]bool)
	count := 0

	for scanner.Scan() {
		var entry struct {
			Word     string `json:"word"`
			LangCode string `json:"lang_code"`
			Pos      string `json:"pos"`
			Senses   []struct {
				Glosses []string `json:"glosses"`
			} `json:"senses"`
		}

		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if entry.LangCode != "fr" || entry.Word == "" {
			continue
		}

		wordLower := strings.ToLower(entry.Word)

		for _, sense := range entry.Senses {
			for _, gloss := range sense.Glosses {
				if gloss == "" {
					continue
				}
				if _, err := stmt.Exec(wordLower, entry.Pos, gloss); err != nil {
					return err
				}

				if !uniqueWords[wordLower] {
					uniqueWords[wordLower] = true
					if _, err := ftsStmt.Exec(wordLower); err != nil {
						return err
					}
				}
				count++
				if count%5000 == 0 {
					fmt.Printf("\rImported: %d entries", count)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("\nImport complete: %d entries (%d unique words)\n", count, len(uniqueWords))
	return nil
}
