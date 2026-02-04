<p align="center">
    <img width="579" height="466" alt="image" src="https://github.com/user-attachments/assets/87e2fa00-1329-4ab0-a311-6c55f203b36f" />
</p>

---

Voca is a CLI language learning & dictionary tool that prioritizes quick access to tailor learning to *your* vocabulary.

## Features

- **Add Words**: Quickly add words with automatic definition lookup.
- **Fuzzy Search**: Interactive, real-time search for your saved collection.
- **Multi-language UI**: Support for different languages in the interface (currently English and French).
- **Language-specific Dictionaries**: Load different dictionary databases based on your target language.
- **Quiz Mode**: Test your vocabulary with interactive quizzes.
- **Import/Export**: Manage your data via CSV or plain text.

## Installation

### From Source

Requirements: Go 1.25+

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/dpardell/voca.git
    cd voca
    ```

2.  **Build and Install**:
    ```bash
    make install
    ```
    By default, this installs the binary to `~/.local/bin` and prepares the data directory at `~/.local/share/voca`. 

    **Customizing Paths**:
    You can use the standard `PREFIX` variable to change the installation root:
    ```bash
    make install PREFIX=/usr/local
    ```
    Or override specific directories:
    ```bash
    make install bindir=/opt/voca/bin datadir=/opt/voca/data
    ```

### Configuration

Voca uses these locations by default (set at build time), but you can override them at runtime using environment variables:

- `VOCA_DB_PATH`: Path to the dictionary database (e.g., `dictionary.db` or `dictionary_fr.db`).
- `VOCA_USER_DB_PATH`: Path to your personal word collection database (default: `$(datadir)/voca.db`).
- `VOCA_LANG`: Set the interface and dictionary language (e.g., `fr`). If unset, it attempts to detect from your `LANG` environment variable.

## Languages

Voca currently supports the following languages:

- **English** (en): Default interface language.
- **French** (fr): Interface and dictionary support.

To switch to French:
```bash
export VOCA_LANG=fr
voca search
```
Only French is supported as an alternative language at this time.

## Usage

### 1. Initialize Dictionary
First, fetch the latest dictionary data for your language (defaults to French if `VOCA_LANG=fr`):
```bash
VOCA_LANG=fr voca install-dict
```

### 2. Add a Word
Interactive mode with fuzzy search:
```bash
voca add
```
Or directly:
```bash
voca add "bonjour"
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
