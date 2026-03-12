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
  - **Interactive Quiz**: Test your knowledge with AI-generated questions tailored to your vocabulary.
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

### Configuration

Voc uses environment variables for runtime configuration. You can also use a `settings.yaml` file in your user config directory.

- **AI Features (Required for Quiz/Convo)** — choose one provider:
  - **Mistral** (simpler setup):
    - `MISTRAL_API_KEY`: Your Mistral API key.
    - `MISTRAL_MODEL`: Model to use (default: `mistral-small-latest`).
  - **Google Vertex AI / Gemini**:
    - `VERTEX_API_KEY`: Your Google Vertex AI/Gemini API key (falls back to `GEMINI_API_KEY`).
    - `VERTEX_PROJECT_ID`: **Required.** Your Google Cloud Project ID.
    - `VERTEX_LOCATION`: Google Cloud Region (default: `us-central1`).
    - `VERTEX_MODEL`: The model ID to use (default: `gemini-2.5-flash-lite`).

  If `MISTRAL_API_KEY` is set, it takes priority over Vertex AI.

- **Application Settings**:
  - `VOC_HOST_LANG`: The language of the UI (e.g., `en`, `fr`).
  - `VOC_TARGET_LANG`: The language you are learning (e.g., `fr`, `es`, `de`).
  - `VOC_DB_PATH`: Path to the dictionary database.
  - `VOC_USER_DB_PATH`: Path to your personal word collection database.

## Languages

Voc supports multiple languages for both the interface and the learning target.

- **Interface (Host Language)**: English (`en`), French (`fr`), Czech (`cs`), Slovak (`sk`), Spanish (`es`), German (`de`), Portuguese (`pt`), and Brazilian Portuguese (`pt-br`) are currently supported.
- **Learning (Target Language)**: Any language available on [Kaikki.org](https://kaikki.org) can be installed (French, Spanish, German, etc.).

To start Voc with specific languages:
```bash
# UI in English, learning French
voc -l en -t fr
```

## Usage

### 1. Interactive Hub
Simply run `voc` to enter the interactive splash screen. From here, you can access all modules.

### 2. Initialize Dictionary
Before searching, install the dictionary for your target language:
```bash
# Installs the dictionary for the current target language
voc install-dict
```
You can also do this directly from the "INSTALL DICTIONARY" option on the splash screen.

### 3. Search & Save
Interactive mode with fuzzy search to find definitions and manage your vocabulary:
```bash
voc search
```

**Keybindings**:
- `Ctrl+P` / `Ctrl+N`: Navigate search results
- `Enter`: View definition
- `Ctrl+S`: Toggle word in your vocabulary (Save/Remove)
- `Esc`: Back / Exit

### 4. Quiz Yourself
Start an AI-powered quiz session tailored to your saved words:
```bash
voc quiz
```

### 5. Practice Conversation
Chat with an AI language coach who will correct your mistakes as you go:
```bash
voc convo
```

### 6. Review Words
List and search through your saved words:
```bash
voc list
```

## Licensing

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Credits

- Dictionary data provided by [kaikki.org](https://kaikki.org).
- TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- CLI framework by [Cobra](https://github.com/spf13/cobra).
