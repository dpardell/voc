# Voc - CLI Language Learning

## Tech Stack
- **Go 1.25+**, Bubble Tea (TUI), Cobra (CLI).
- **SQLite3** with `fts5` support (required).
- **Vertex AI** via `cloud.google.com/go/vertexai/genai`.

## Architecture
- `cmd/voc/`: CLI commands.
- `internal/database/`: SQLite management for personal vocab and progress.
- `internal/dictionary/`: Kaikki.org data ingestion and FTS5 search.
- `internal/llm/`: Prompts and Vertex AI integration.
- `internal/ui/`: Bubble Tea models and central styling.
- `internal/i18n/`: Multi-language UI logic.

## Build & Run
- **Build**: `go build -tags "fts5" ./cmd/voc` (or `make build`).
- **Install**: `make install` (defaults to `~/.local/bin` and `~/.local/share/voc`).
- **Test**: `go test ./...`.

## Configuration
- `VERTEX_API_KEY`: API key (falls back to `GEMINI_API_KEY`).
- `VERTEX_PROJECT_ID`: **Required** Google Cloud Project ID.
- `VERTEX_MODEL`: Default is `gemini-2.5-flash-lite`.
- `VOC_HOST_LANG` / `VOC_TARGET_LANG`: Language settings.

## Development
- Logic stays in `internal/`.
- Use `sanitizeJSON` helper for LLM outputs.
- Adhere to `internal/ui/styles.go` for TUI consistency.
