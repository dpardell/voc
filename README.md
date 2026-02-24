<p align="center">
    <img width="400" height="411" alt="image" src="https://github.com/user-attachments/assets/ad10b029-7533-4931-8c0b-7aafcf9f95b0" />
</p>

---

Voc is a CLI language learning & dictionary tool that prioritizes quick access to tailor learning to *your* vocabulary.

***A note from David:**
This project started when I was a software developer in a French class. I wanted a way to record new words as they came to me with as little friction as possible.
That tool was originally written in Ruby.
I have since put a some effort into cleaning this up and rewritting it in Go, but it is still very much a tool I am building for myself. That being said, contributions
are defintely welcome. Long term, I would like for this project to head in the direction of a more complete CLI-based langugage learning app.*

## Features

- **Dictionary Search**: Interactive, real-time search in local dictionary with text indexing.
- **Personal Vocabulary**: Save words to your personal vocabulary.
- **AI-Powered Learning**:
  - **Interactive Quiz**: Test your knowledge with a mix of multiple-choice and fill-in-the-blank questions.
  - **Language Coach (`voc convo`)**: Chat with a friendly AI coach that provides real-time grammar and stylistic corrections.
  - **Progress Tracking**: A persistent `progress.md` file keeps a high-level overview of your grammar competency and vocabulary range.
- **Multi-language UI**: Support for different languages in the interface (currently English and French).
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

- `GEMINI_API_KEY`: **Required for AI features.** Your Google Gemini API key.
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

### 2. Search & Save
Interactive mode with fuzzy search to find definitions and manage your vocabulary:
```bash
voc search
```

**Keybindings**:
Search results:
- `Ctrl+P` / `Ctrl+N`: Navigate search results
- `Enter`: View definition

Definition view:
- `J` / `K`: Scroll definition
- `Ctrl+S`: Toggle word in your vocabulary (Save/Remove)
- `Esc`: Back / Exit

### 3. Quiz Yourself
Start an AI-powered quiz session tailored to your current progress:
```bash
voc quiz
```
*Note: Use `voc quiz --flashcards` for the traditional offline flashcard mode.*

### 4. Practice Conversation
Chat with an AI language coach who will correct your mistakes as you go:
```bash
voc convo
```
This mode uses a "messaging app" style interface with speech bubbles and a dedicated area for grammatical feedback.

### 5. Review Words
List all saved words:
```bash
voc list
```

## Licensing

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

- Dictionary data provided by [kaikki.org](https://kaikki.org).
- TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- CLI framework by [Cobra](https://github.com/spf13/cobra).
