# Voca

Voca is a personal CLI dictionary tool for language learners. It helps you track new words, automatically fetches definitions from Wiktionary, and provides spaced-repetition style quizzes.

## Features

- **Add Words**: Quickly add words with automatic definition lookup.
- **Fuzzy Search**: Interactive, real-time search for your saved collection.
- **Dictionary**: Integrated with the French Wiktionary (via [kaikki.org](https://kaikki.org)).
- **Quiz Mode**: Test your vocabulary with interactive quizzes.
- **Import/Export**: Manage your data via CSV or plain text.
- **No Dependencies**: Single binary executable (SQLite embedded).

## Installation

### From Source

Requirements: Go 1.21+

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/david/voca.git
    cd voca
    ```

2.  **Build and Install**:
    ```bash
    make install
    ```
    This will install the `voca` binary to `~/.local/bin`. Ensure this directory is in your `$PATH`.

## Usage

### 1. Initialize Dictionary
First, fetch the latest French dictionary data:
```bash
voca install-dict
```

### 2. Add a Word
Interactive mode with fuzzy search:
```bash
voca add
```
Or directly:
```bash
voca add "bonjour" -c "Hello in French"
```

### 3. Review Words
List all words:
```bash
voca list
```
Show specific details:
```bash
voca show "bonjour"
```

### 4. Quiz Yourself
Start a quiz session:
```bash
voca quiz
```

## Licensing

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

- Dictionary data provided by [kaikki.org](https://kaikki.org).
- TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- CLI framework by [Cobra](https://github.com/spf13/cobra).
