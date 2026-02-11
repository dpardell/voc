package i18n

import (
	"fmt"
	"os"
)

type StringID string

const (
	NoMatches              StringID = "NoMatches"
	SelectItem             StringID = "SelectItem"
	ErrorPreview           StringID = "ErrorPreview"
	ResizeWindow           StringID = "ResizeWindow"
	TypeToSearch           StringID = "TypeToSearch"
	Searching              StringID = "Searching"
	AddingWord             StringID = "AddingWord"
	FixingWord             StringID = "FixingWord"
	ShowingWord            StringID = "ShowingWord"
	DictionaryNotInstalled StringID = "DictionaryNotInstalled"
	ErrWordNotFound        StringID = "ErrWordNotFound"
	ErrDatabase            StringID = "ErrDatabase"
	DictionaryNotFound     StringID = "DictionaryNotFound"
	PromptDownloadDefault  StringID = "PromptDownloadDefault"
	PromptCustomURL        StringID = "PromptCustomURL"
	InstallCancelled       StringID = "InstallCancelled"
	HintsSearch            StringID = "HintsSearch"
	HintsDefinition        StringID = "HintsDefinition"
	Saved                  StringID = "Saved"
	NotSaved               StringID = "NotSaved"
	WordCount              StringID = "WordCount"
)

var en = map[StringID]string{
	NoMatches:              "No matches found...",
	SelectItem:             "Select an item to see details...",
	ErrorPreview:           "Error loading preview",
	ResizeWindow:           "(Resize window to view results)",
	TypeToSearch:           "Type to search...",
	Searching:              "SEARCH DICTIONARY",
	AddingWord:             "ADD WORD",
	FixingWord:             "FIXING WORD: %s",
	ShowingWord:            "SHOW WORD",
	DictionaryNotInstalled: "Dictionary not installed. Run 'voc install-dict' first.",
	ErrWordNotFound:        "Word not found in dictionary",
	ErrDatabase:            "Error initializing database",
	DictionaryNotFound:     "Dictionary not found at %s",
	PromptDownloadDefault:  "Download dictionary from %s? [Y/n/c (custom)] ",
	PromptCustomURL:        "Enter custom dictionary URL: ",
	InstallCancelled:       "Installation cancelled.",
	HintsSearch:            "Ctrl-C: Quit • Ctrl-P: Up • Enter: Select",
	HintsDefinition:        "Esc: Back • J/K: Scroll • Ctrl-P/N: Prev/Next • Ctrl-S: Toggle Vocab",
	Saved:                  "SAVED ❤️",
	NotSaved:               "NOT SAVED",
	WordCount:              "Your List (%d words)",
}

var fr = map[StringID]string{
	NoMatches:              "Aucune correspondance trouvée...",
	SelectItem:             "Sélectionnez un élément pour voir les détails...",
	ErrorPreview:           "Erreur lors du chargement de l'aperçu",
	ResizeWindow:           "(Redimensionnez la fenêtre pour afficher les résultats)",
	TypeToSearch:           "Tapez pour rechercher...",
	Searching:              "RECHERCHER DANS LE DICTIONNAIRE",
	AddingWord:             "AJOUTER UN MOT",
	FixingWord:             "CORRIGER LE MOT : %s",
	ShowingWord:            "AFFICHER LE MOT",
	DictionaryNotInstalled: "Dictionnaire non installé. Exécutez 'voc install-dict' d'abord.",
	ErrWordNotFound:        "Mot non trouvé dans le dictionnaire",
	ErrDatabase:            "Erreur lors de l'initialisation de la base de données",
	DictionaryNotFound:     "Dictionnaire non trouvé à %s",
	PromptDownloadDefault:  "Télécharger le dictionnaire depuis %s ? [Y/n/c (personnalisé)] ",
	PromptCustomURL:        "Entrez l'URL personnalisée du dictionnaire : ",
	InstallCancelled:       "Installation annulée.",
	HintsSearch:            "Ctrl-C: Quitter • Ctrl-P: Haut • Entrée: Sélectionner",
	HintsDefinition:        "Esc: Retour • J/K: Défiler • Ctrl-P/N: Préc/Suiv • Ctrl-S: Vocabulaire",
	Saved:                  "ENREGISTRÉ ❤️",
	NotSaved:               "NON ENREGISTRÉ",
	WordCount:              "Votre Liste (%d mots)",
}

var activeLocale = "en"

var locales = map[string]map[StringID]string{
	"en": en,
	"fr": fr,
}

func init() {
	lang := os.Getenv("VOC_LANG")
	if lang == "" {
		lang = os.Getenv("LANG")
	}
	if len(lang) >= 2 {
		SetLanguage(lang[:2])
	}
}

func SetLanguage(lang string) {
	if _, ok := locales[lang]; ok {
		activeLocale = lang
	}
}

func GetLanguage() string {
	return activeLocale
}

func T(id StringID, args ...any) string {
	translations, ok := locales[activeLocale]
	if !ok {
		translations = en
	}

	val, ok := translations[id]
	if !ok {
		return string(id)
	}

	if len(args) > 0 {
		return fmt.Sprintf(val, args...)
	}
	return val
}
