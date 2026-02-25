package dictionary

import (
	"bufio"
	"compress/gzip"
	"context"
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

var KaikkiURL string

type Importer struct {
	DictDir   string
	DictDB    string
	CacheFile string
	Lang      string
}

func NewImporter(lang string) (*Importer, error) {
	dbPath, err := GetDictionaryPath(lang)
	if err != nil {
		return nil, err
	}

	dictDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dictDir, 0755); err != nil {
		return nil, err
	}

	return &Importer{
		DictDir:   dictDir,
		DictDB:    dbPath,
		CacheFile: filepath.Join(dictDir, fmt.Sprintf("kaikki-%s.jsonl", lang)),
		Lang:      lang,
	}, nil
}

func (i *Importer) DownloadAndImport(ctx context.Context, force bool, url string) error {
	if _, err := os.Stat(i.CacheFile); err == nil && !force {
		fmt.Println("Using cached dictionary file...")
		fmt.Println("  (use --force to re-download)")
		return i.importFile(ctx)
	}

	fmt.Printf("Downloading %s dictionary from %s...\n", i.Lang, url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Use a buffered reader to peek at the magic bytes
	reader := bufio.NewReader(resp.Body)
	isGzip := false
	peek, err := reader.Peek(2)
	if err == nil && peek[0] == 0x1f && peek[1] == 0x8b {
		isGzip = true
	}

	var finalReader io.Reader = reader
	if isGzip {
		fmt.Println("Decompressing and saving...")
		gzReader, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer gzReader.Close()
		finalReader = gzReader
	} else {
		fmt.Println("Saving...")
	}

	outFile, err := os.Create(i.CacheFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, finalReader)
	if err != nil {
		return err
	}

	return i.importFile(ctx)
}

func (i *Importer) importFile(ctx context.Context) error {
	fmt.Println("Importing into database...")

	tempDBPath := i.DictDB + ".tmp"
	os.Remove(tempDBPath)

	db, err := sql.Open("sqlite3", tempDBPath)
	if err != nil {
		return err
	}
	defer func() {
		db.Close()
		os.Remove(tempDBPath) // Cleanup if we fail before rename
	}()

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
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

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

		if entry.LangCode != i.Lang || entry.Word == "" {
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

	db.Close() // Close before rename

	if err := os.Rename(tempDBPath, i.DictDB); err != nil {
		return fmt.Errorf("failed to finalize database: %v", err)
	}

	fmt.Printf("\nImport complete: %d entries (%d unique words)\n", count, len(uniqueWords))
	return nil
}
