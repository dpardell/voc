# Voc - CLI Language Learning & Dictionary

Voc is a CLI-based language learning tool written in Go. it combines traditional dictionary functionality with AI-powered features like interactive quizzes and a conversation coach.

## Tech Stack
- **Language**: Go 1.25.1+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) & [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **Database**: SQLite3 (requires `fts5` support)
- **AI Integration**: [Google Vertex AI API](https://cloud.google.com/vertex-ai) via `cloud.google.com/go/vertexai/genai`

## Architecture
The project follows a standard Go CLI structure:
- `cmd/voc/`: Contains the CLI command definitions (Cobra). Each subcommand (`search`, `quiz`, `convo`, etc.) has its own file.
- `internal/database/`: Manages the local SQLite database for personal vocabulary.
- `internal/dictionary/`: Handles dictionary data ingestion (from Kaikki.org) and searching.
- `internal/llm/`: Encapsulates AI logic, including prompt engineering for quizzes and the conversation coach.
- `internal/ui/`: Contains TUI components, models, and centralized styling.
- `internal/i18n/`: Logic for multi-language interface support (currently English and French).

## Building and Running

### Key Commands
- **Build**: `make build` (Produces the `voc` binary in the root directory).
- **Install**: `make install` (Installs to `~/.local/bin` and sets up data in `~/.local/share/voc` by default).
- **Test**: `go test ./...`
- **Clean**: `make clean`

### Build Requirements
The project uses SQLite with Full Text Search 5. When building manually, ensure the `fts5` tag is included:
```bash
go build -tags "fts5" ./cmd/voc
```

### Environment Variables
- `VERTEX_API_KEY`: Required for AI-powered features (`quiz`, `convo`). Falls back to `GEMINI_API_KEY` if not set.
- `VERTEX_PROJECT_ID`: **Required**. Your Google Cloud Project ID.
- `VERTEX_LOCATION`: Optional. Your Google Cloud Region (defaults to `us-central1`).
- `VERTEX_MODEL`: Optional. The model ID to use (defaults to `gemini-2.5-flash-lite`).
- `VOC_LANG`: Sets the interface and dictionary language (e.g., `fr`).
- `VOC_DB_PATH`: Path to the dictionary database.
- `VOC_USER_DB_PATH`: Path to the personal vocabulary database.

## Development Conventions

### Coding Style
- **Internal Packages**: All core logic resides in `internal/` to prevent external imports and maintain a clean API.
- **TUI Models**: Bubble Tea models are used for interactive components. Check `internal/ui/` for existing patterns.
- **Styling**: TUI colors and styles are centralized in `internal/ui/styles.go`. Adhere to these styles for a consistent look.

### AI Integration
- LLM prompts are located in `internal/llm/llm.go`. 
- The project targets `gemini-2.5-flash-lite` (configurable) on Google Cloud Vertex AI.
- Always use the `sanitizeJSON` helper when expecting structured JSON output from the LLM.

### Testing
- Tests are located alongside source files (e.g., `llm_test.go`).
- Use `go test ./...` to run all tests.
- When adding new features, include corresponding tests in the `internal/` packages.

## Data Sources
- Dictionary data is fetched from [kaikki.org](https://kaikki.org) and imported into a local SQLite database with FTS5 enabled for fast searching.
