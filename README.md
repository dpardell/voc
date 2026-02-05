<p align="center">
    <img width="400" height="411" alt="image" src="https://github.com/user-attachments/assets/ad10b029-7533-4931-8c0b-7aafcf9f95b0" />
</p>

---

Voc is a CLI language learning & dictionary tool that prioritizes quick access to tailor learning to *your* vocabulary.

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
    git clone https://github.com/dpardell/voc.git
    cd voc
    ```

2.  **Build and Install**:
    ```bash
    make install
    ```
    By default, this installs the binary to `~/.local/bin` and prepares the data directory at `~/.local/share/voc`. 

    **Customizing Paths**:
    You can use the standard `PREFIX` variable to change the installation root:
    ```bash
    make install PREFIX=/usr/local
    ```
    Or override specific directories:
    ```bash
    make install bindir=/opt/voc/bin datadir=/opt/voc/data
    ```

### Configuration

Voc uses these locations by default (set at build time), but you can override them at runtime using environment variables:

- `VOC_DB_PATH`: Path to the dictionary database (e.g., `dictionary.db` or `dictionary_fr.db`).
- `VOC_USER_DB_PATH`: Path to your personal word collection database (default: `$(datadir)/voc.db`).
- `VOC_LANG`: Set the interface and dictionary language (e.g., `fr`). If unset, it attempts to detect from your `LANG` environment variable.

## Languages

Voc currently supports the following languages:

- **English** (en): Default interface language.
- **French** (fr): Interface and dictionary support.

To switch to French:
```bash
export VOC_LANG=fr
voc search
```
Only French is supported as an alternative language at this time.

## Usage

### 1. Initialize Dictionary
First, fetch the latest dictionary data for your language (defaults to French if `VOC_LANG=fr`):
```bash
VOC_LANG=fr voc install-dict
```

### 2. Add a Word
Interactive mode with fuzzy search:
```bash
voc add
```
Or directly:
```bash
voc add "bonjour"
```

### 3. Review Words
List all words:
```bash
voc list
```
Show specific details:
```bash
voc show "bonjour"
```

### 4. Quiz Yourself
Start a quiz session:
```bash
voc quiz
```

## Licensing

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

- Dictionary data provided by [kaikki.org](https://kaikki.org).
- TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- CLI framework by [Cobra](https://github.com/spf13/cobra).
