# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build (fts5 tag required for SQLite full-text search)
make build

# Run tests
make test

# Run a single test package
go test -tags "fts5" ./internal/config/...

# Run a single test by name
go test -tags "fts5" -run TestFunctionName ./internal/config/...

# Install to ~/.local/bin
make install

# Install to custom prefix
make install PREFIX=/usr/local
```

> **Important:** Always include `-tags "fts5"` when running `go build` or `go test` directly. The Makefile handles this automatically.

## Architecture

The entry point is `cmd/voc/main.go`, which wires a Cobra CLI with subcommands (`search`, `quiz`, `convo`, `list`, `install-dict`). Running `voc` with no arguments launches the interactive splash screen TUI, which is a loop that delegates back into the same subcommand handlers.

`internal/app/app.go` is the central dependency container passed to all UI components. It holds:
- `DB` — the user's personal vocabulary (SQLite via `go-sqlite3`)
- `Dict` — the read-only dictionary (separate SQLite, FTS5-indexed, sourced from kaikki.org)
- `HostLang` / `TargetLang` — resolved at startup
- `Settings` — config loaded from file + env

### Package overview

| Package | Role |
|---|---|
| `internal/app` | Creates and owns all shared dependencies; provides `GetLLMClient()` |
| `internal/config` | Loads/saves `settings.yaml`; config precedence: defaults → file → env vars → CLI flags |
| `internal/database` | User vocabulary SQLite DB (`words`, `word_types`, `definitions` tables + schema migrations); also reads/writes `progress.md` |
| `internal/dictionary` | Read-only dictionary SQLite (FTS5 `words_fts` table with LIKE fallback); installed per target language as `dictionary_{lang}.db` |
| `internal/llm` | Google Vertex AI (Gemini) client — `GenerateQuiz`, `Chat`, `UpdateProgress` |
| `internal/i18n` | Static string maps per locale (`en`, `fr`, `cs`, `sk`, `es`, `de`, `pt`, `pt-br`); `T(StringID)` for lookups; `SetHostLanguage` / `SetTargetLanguage` set package-level globals |
| `internal/ui` | Bubble Tea TUI models: `splash`, `fuzzy` (search), `quiz`, `convo` |

### Data flow

1. `config.Load()` reads `~/.config/voc/settings.yaml`, then overrides with `VOC_HOST_LANG` / `VOC_TARGET_LANG` env vars, then CLI flags.
2. The dictionary DB path resolves as: `VOC_DB_PATH` env → `DefaultDictionaryDirectory` build var (set by `make`) → `~/.config/voc/dictionary_{lang}.db`.
3. The user DB path resolves as: `VOC_USER_DB_PATH` env → `DefaultUserDBPath` build var → `~/.local/share/voc/voc.db`.
4. Both `DefaultDictionaryDirectory` and `DefaultUserDBPath` are injected at build time via `-ldflags` in the Makefile.

### Adding a new UI language

Add a new `map[StringID]string` in `internal/i18n/i18n.go` covering every `StringID` constant, then register it in the `locales` map and `GetLanguageName`.

### TUI styling

All Lipgloss styles are defined in `internal/ui/styles.go`. Use these when building or modifying UI components rather than defining styles inline.

### LLM integration

`internal/llm` wraps `cloud.google.com/go/vertexai/genai`. Requires env vars `VERTEX_API_KEY` (or `GEMINI_API_KEY`) and `VERTEX_PROJECT_ID`. Default model is `gemini-2.5-flash-lite`. All LLM responses are JSON; `sanitizeJSON` strips markdown code fences before unmarshalling.
