package i18n

import (
	"fmt"
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
	QuizTitle              StringID = "QuizTitle"
	QuizQuestion           StringID = "QuizQuestion"
	QuizScore              StringID = "QuizScore"
	QuizCorrect            StringID = "QuizCorrect"
	QuizWrong              StringID = "QuizWrong"
	QuizCompleted          StringID = "QuizCompleted"
	HintsQuiz              StringID = "HintsQuiz"
	// New Strings
	ConvoPlaceholder  StringID = "ConvoPlaceholder"
	ConvoStarted      StringID = "ConvoStarted"
	ConvoCorrections  StringID = "ConvoCorrections"
	ConvoThinking     StringID = "ConvoThinking"
	ConvoError        StringID = "ConvoError"
	ConvoExitHint     StringID = "ConvoExitHint"
	ConvoCoachTitle   StringID = "ConvoCoachTitle"
	QuizPlaceholder   StringID = "QuizPlaceholder"
	PressEnterToCont  StringID = "PressEnterToCont"
	ErrorNoWords      StringID = "ErrorNoWords"
	ErrorNoSavedWords StringID = "ErrorNoSavedWords"
	DictionaryLoading StringID = "DictionaryLoading"
	// Splash screen
	SplashTagline StringID = "SplashTagline"
	MenuSearch    StringID = "MenuSearch"
	MenuQuiz      StringID = "MenuQuiz"
	MenuConvo     StringID = "MenuConvo"
	MenuVocab     StringID = "MenuVocab"
	MenuInstall   StringID = "MenuInstall"
	MenuQuit      StringID = "MenuQuit"
	SplashFooter  StringID = "SplashFooter"
	// Import/Export
	ImportComplete StringID = "ImportComplete"
	ImportUsage    StringID = "ImportUsage"
	ExportSuccess  StringID = "ExportSuccess"
	// Quiz states
	QuizGenerating StringID = "QuizGenerating"
	QuizUpdating   StringID = "QuizUpdating"
	QuizHintType   StringID = "QuizHintType"
	// CLI extra strings
	FlashcardsWelcome StringID = "FlashcardsWelcome"
	FlashcardsPress   StringID = "FlashcardsPress"
	FlashcardsNext    StringID = "FlashcardsNext"
	FlashcardsEnd     StringID = "FlashcardsEnd"
	DailySayingPrompt StringID = "DailySayingPrompt"
	NoWordsToQuiz     StringID = "NoWordsToQuiz"
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
	QuizTitle:              "VOCABULARY QUIZ",
	QuizQuestion:           "Question %d/%d",
	QuizScore:              "Score: %d/%d",
	QuizCorrect:            "CORRECT! ✨",
	QuizWrong:              "WRONG! The correct answer was: %s",
	QuizCompleted:          "Quiz completed! Final Score: %d/%d",
	HintsQuiz:              "Enter: Submit/Next • Ctrl-C: Quit",
	ConvoPlaceholder:       "Say something in your target language...",
	ConvoStarted:           "Chat started! Say hi!",
	ConvoCorrections:       "Corrections:",
	ConvoThinking:          "Thinking...",
	ConvoError:             "ERROR",
	ConvoExitHint:          "Press Ctrl+C to exit",
	ConvoCoachTitle:        "LANGUAGE COACH",
	QuizPlaceholder:        "Type your answer...",
	PressEnterToCont:       "Press Enter to continue...",
	ErrorNoWords:           "Not enough words to generate a quiz. Add more words first!",
	ErrorNoSavedWords:      "You haven't saved any words yet. Search for words to add them!",
	DictionaryLoading:      "Dictionary not loaded",
	SplashTagline:          "Your AI-Powered Language Companion",
	MenuSearch:             "SEARCH DICTIONARY",
	MenuQuiz:               "VOCABULARY QUIZ",
	MenuConvo:              "LANGUAGE COACH",
		MenuVocab:     "MY VOCABULARY",
		MenuInstall:   "INSTALL DICTIONARY",
		MenuQuit:      "QUIT",
	SplashFooter:           "↑/↓: navigate • enter: select • q: quit",
	ImportComplete:         "Import complete: %d added, %d skipped",
	ImportUsage:            "Usage: voc import -f <file> OR pipe input",
	ExportSuccess:          "Exported %d words to %s",
	QuizGenerating:         "Generating quiz...",
	QuizUpdating:           "Updating progress...",
	QuizHintType:           "Type your answer and press Enter",
	FlashcardsWelcome:      "Flashcards mode - Press Enter to see definition, Ctrl+C to exit",
	FlashcardsPress:        "Press Enter to see definition... ",
	FlashcardsNext:         "Press Enter for next word... ",
	FlashcardsEnd:          "Quiz ended.",
	DailySayingPrompt:      "Welcome back! (Set VERTEX_API_KEY/PROJECT_ID for daily sayings)",
	NoWordsToQuiz:          "No words to quiz.",
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
	QuizTitle:              "QUIZ DE VOCABULAIRE",
	QuizQuestion:           "Question %d/%d",
	QuizScore:              "Score : %d/%d",
	QuizCorrect:            "CORRECT ! ✨",
	QuizWrong:              "MAUVAIS ! La bonne réponse était : %s",
	QuizCompleted:          "Quiz terminé ! Score final : %d/%d",
	HintsQuiz:              "Entrée : Valider/Suivant • Ctrl-C : Quitter",
	ConvoPlaceholder:       "Dites quelque chose dans votre langue cible...",
	ConvoStarted:           "Chat commencé ! Dites bonjour !",
	ConvoCorrections:       "Corrections :",
	ConvoThinking:          "Réflexion...",
	ConvoError:             "ERREUR",
	ConvoExitHint:          "Appuyez sur Ctrl+C pour quitter",
	ConvoCoachTitle:        "COACH LINGUISTIQUE",
	QuizPlaceholder:        "Tapez votre réponse...",
	PressEnterToCont:       "Appuyez sur Entrée pour continuer...",
	ErrorNoWords:           "Pas assez de mots pour générer un quiz. Ajoutez plus de mots d'abord !",
	ErrorNoSavedWords:      "Vous n'avez pas encore enregistré de mots. Recherchez des mots pour les ajouter !",
	DictionaryLoading:      "Dictionnaire non chargé",
	SplashTagline:          "Votre compagnon linguistique propulsé par l'IA",
	MenuSearch:             "RECHERCHER DANS LE DICTIONNAIRE",
	MenuQuiz:               "QUIZ DE VOCABULAIRE",
	MenuConvo:              "COACH LINGUISTIQUE",
		MenuVocab:     "MON VOCABULAIRE",
		MenuInstall:   "INSTALLER LE DICTIONNAIRE",
		MenuQuit:      "QUITTER",
	SplashFooter:           "↑/↓: naviguer • entrée: sélectionner • q: quitter",
	ImportComplete:         "Importation terminée : %d ajouté(s), %d ignoré(s)",
	ImportUsage:            "Utilisation : voc import -f <fichier> OU pipe entrée",
	ExportSuccess:          "Exportation de %d mots vers %s",
	QuizGenerating:         "Génération du quiz...",
	QuizUpdating:           "Mise à jour de la progression...",
	QuizHintType:           "Tapez votre réponse et appuyez sur Entrée",
	FlashcardsWelcome:      "Mode Flashcards - Appuyez sur Entrée pour voir la définition, Ctrl+C pour quitter",
	FlashcardsPress:        "Appuyez sur Entrée pour voir la définition... ",
	FlashcardsNext:         "Appuyez sur Entrée pour le mot suivant... ",
	FlashcardsEnd:          "Quiz terminé.",
	DailySayingPrompt:      "Bon retour ! (Configurez VERTEX_API_KEY/PROJECT_ID pour les dictons quotidiens)",
	NoWordsToQuiz:          "Aucun mot pour le quiz.",
}

var cs = map[StringID]string{
	NoMatches:              "Žádné shody nenalezeny...",
	SelectItem:             "Vyberte položku pro zobrazení detailů...",
	ErrorPreview:           "Chyba při načítání náhledu",
	ResizeWindow:           "(Změňte velikost okna pro zobrazení výsledků)",
	TypeToSearch:           "Pište pro hledání...",
	Searching:              "HLEDAT VE SLOVNÍKU",
	AddingWord:             "PŘIDAT SLOVO",
	FixingWord:             "OPRAVA SLOVA: %s",
	ShowingWord:            "ZOBRAZIT SLOVO",
	DictionaryNotInstalled: "Slovník není nainstalován. Spusťte nejdříve 'voc install-dict'.",
	ErrWordNotFound:        "Slovo nebylo nalezeno ve slovníku",
	ErrDatabase:            "Chyba při inicializaci databáze",
	DictionaryNotFound:     "Slovník nebyl nalezen v %s",
	PromptDownloadDefault:  "Stáhnout slovník z %s? [Y/n/c (vlastní)] ",
	PromptCustomURL:        "Zadejte vlastní URL slovníku: ",
	InstallCancelled:       "Instalace zrušena.",
	HintsSearch:            "Ctrl-C: Konec • Ctrl-P: Nahoru • Enter: Vybrat",
	HintsDefinition:        "Esc: Zpět • J/K: Posun • Ctrl-P/N: Předchozí/Další • Ctrl-S: Přepnout slovní zásobu",
	Saved:                  "ULOŽENO ❤️",
	NotSaved:               "NEULOŽENO",
	WordCount:              "Váš seznam (%d slov)",
	QuizTitle:              "KVÍZ ZE SLOVNÍ ZÁSOBY",
	QuizQuestion:           "Otázka %d/%d",
	QuizScore:              "Skóre: %d/%d",
	QuizCorrect:            "SPRÁVNĚ! ✨",
	QuizWrong:              "ŠPATNĚ! Správná odpověď byla: %s",
	QuizCompleted:          "Kvíz dokončen! Konečné skóre: %d/%d",
	HintsQuiz:              "Enter: Odeslat/Další • Ctrl-C: Konec",
	ConvoPlaceholder:       "Řekněte něco ve vašem cílovém jazyce...",
	ConvoStarted:           "Chat zahájen! Pozdravte!",
	ConvoCorrections:       "Opravy:",
	ConvoThinking:          "Přemýšlím...",
	ConvoError:             "CHYBA",
	ConvoExitHint:          "Stiskněte Ctrl+C pro ukončení",
	ConvoCoachTitle:        "JAZYKOVÝ KOUČ",
	QuizPlaceholder:        "Napište svou odpověď...",
	PressEnterToCont:       "Stiskněte Enter pro pokračování...",
	ErrorNoWords:           "Nedostatek slov pro vygenerování kvízu. Nejdříve přidejte slova!",
	ErrorNoSavedWords:      "Zatím jste neuložili žádná slova. Hledejte slova a přidejte je!",
	DictionaryLoading:      "Slovník není načten",
	SplashTagline:          "Váš jazykový společník s umělou inteligencí",
	MenuSearch:             "HLEDAT VE SLOVNÍKU",
	MenuQuiz:               "KVÍZ ZE SLOVNÍ ZÁSOBY",
	MenuConvo:              "JAZYKOVÝ KOUČ",
	MenuVocab:              "MOJE SLOVNÍ ZÁSOBA",
	MenuInstall:            "INSTALOVAT SLOVNÍK",
	MenuQuit:               "KONEC",
	SplashFooter:           "↑/↓: navigace • enter: vybrat • q: konec",
	ImportComplete:         "Import dokončen: %d přidáno, %d přeskočeno",
	ImportUsage:            "Použití: voc import -f <soubor> NEBO pipe vstup",
	ExportSuccess:          "Exportováno %d slov do %s",
	QuizGenerating:         "Generování kvízu...",
	QuizUpdating:           "Aktualizace pokroku...",
	QuizHintType:           "Napište svou odpověď a stiskněte Enter",
	FlashcardsWelcome:      "Režim kartiček - Stiskněte Enter pro zobrazení definice, Ctrl+C pro ukončení",
	FlashcardsPress:        "Stiskněte Enter pro zobrazení definice... ",
	FlashcardsNext:         "Stiskněte Enter pro další slovo... ",
	FlashcardsEnd:          "Kvíz ukončen.",
	DailySayingPrompt:      "Vítejte zpět! (Nastavte VERTEX_API_KEY/PROJECT_ID pro denní rčení)",
	NoWordsToQuiz:          "Žádná slova pro kvíz.",
}

var sk = map[StringID]string{
	NoMatches:              "Žiadne zhody sa nenašli...",
	SelectItem:             "Vyberte položku pre zobrazenie detailov...",
	ErrorPreview:           "Chyba pri načítaní náhľadu",
	ResizeWindow:           "(Zmeňte veľkosť okna pre zobrazenie výsledkov)",
	TypeToSearch:           "Píšte pre hľadanie...",
	Searching:              "HĽADAŤ VO SLOVNÍKU",
	AddingWord:             "PRIDAŤ SLOVO",
	FixingWord:             "OPRAVA SLOVA: %s",
	ShowingWord:            "ZOBRAZIŤ SLOVO",
	DictionaryNotInstalled: "Slovník nie je nainštalovaný. Spustite najskôr 'voc install-dict'.",
	ErrWordNotFound:        "Slovo sa nenašlo v slovníku",
	ErrDatabase:            "Chyba pri inicializácii databázy",
	DictionaryNotFound:     "Slovník sa nenašiel v %s",
	PromptDownloadDefault:  "Stiahnuť slovník z %s? [Y/n/c (vlastný)] ",
	PromptCustomURL:        "Zadajte vlastnú URL slovníka: ",
	InstallCancelled:       "Inštalácia zrušená.",
	HintsSearch:            "Ctrl-C: Koniec • Ctrl-P: Hore • Enter: Vybrať",
	HintsDefinition:        "Esc: Späť • J/K: Posun • Ctrl-P/N: Predchádzajúce/Ďalšie • Ctrl-S: Prepnúť slovnú zásobu",
	Saved:                  "ULOŽENÉ ❤️",
	NotSaved:               "NEULOŽENÉ",
	WordCount:              "Váš zoznam (%d slov)",
	QuizTitle:              "KVÍZ ZO SLOVNEJ ZÁSOBY",
	QuizQuestion:           "Otázka %d/%d",
	QuizScore:              "Skóre: %d/%d",
	QuizCorrect:            "SPRÁVNE! ✨",
	QuizWrong:              "NESPRÁVNE! Správna odpoveď bola: %s",
	QuizCompleted:          "Kvíz dokončený! Konečné skóre: %d/%d",
	HintsQuiz:              "Enter: Odoslať/Ďalej • Ctrl-C: Koniec",
	ConvoPlaceholder:       "Povedzte niečo vo vašom cieľovom jazyku...",
	ConvoStarted:           "Chat zahájený! Pozdravte!",
	ConvoCorrections:       "Opravy:",
	ConvoThinking:          "Premýšľam...",
	ConvoError:             "CHYBA",
	ConvoExitHint:          "Stlačte Ctrl+C pre ukončenie",
	ConvoCoachTitle:        "JAZYKOVÝ KOUČ",
	QuizPlaceholder:        "Napíšte svoju odpoveď...",
	PressEnterToCont:       "Stlačte Enter pre pokračovanie...",
	ErrorNoWords:           "Nedostatok slov na vygenerovanie kvízu. Najskôr pridajte slová!",
	ErrorNoSavedWords:      "Zatiaľ ste neuložili žiadne slová. Hľadajte slová a pridajte ich!",
	DictionaryLoading:      "Slovník nie je načítaný",
	SplashTagline:          "Váš jazykový spoločník s umelou inteligenciou",
	MenuSearch:             "HĽADAŤ VO SLOVNÍKU",
	MenuQuiz:               "KVÍZ ZO SLOVNEJ ZÁSOBY",
	MenuConvo:              "JAZYKOVÝ KOUČ",
	MenuVocab:              "MOJA SLOVNÁ ZÁSOBA",
	MenuInstall:            "INŠTALOVAŤ SLOVNÍK",
	MenuQuit:               "KONIEC",
	SplashFooter:           "↑/↓: navigácia • enter: vybrať • q: koniec",
	ImportComplete:         "Import dokončený: %d pridaných, %d preskočených",
	ImportUsage:            "Použitie: voc import -f <súbor> ALEBO pipe vstup",
	ExportSuccess:          "Exportovaných %d slov do %s",
	QuizGenerating:         "Generovanie kvízu...",
	QuizUpdating:           "Aktualizácia pokroku...",
	QuizHintType:           "Napíšte svoju odpoveď a stlačte Enter",
	FlashcardsWelcome:      "Režim kartičiek - Stlačte Enter pre zobrazenie definície, Ctrl+C pre ukončenie",
	FlashcardsPress:        "Stlačte Enter pre zobrazenie definice... ",
	FlashcardsNext:         "Stlačte Enter pre ďalšie slovo... ",
	FlashcardsEnd:          "Kvíz ukončený.",
	DailySayingPrompt:      "Vitajte späť! (Nastavte VERTEX_API_KEY/PROJECT_ID pre denné porekadlá)",
	NoWordsToQuiz:          "Žiadne slová pre kvíz.",
}

var activeHostLang = "en"
var activeTargetLang = "fr"

var locales = map[string]map[StringID]string{
	"en": en,
	"fr": fr,
	"cs": cs,
	"sk": sk,
}

func init() {
	// Defaults will be overridden by config loading in internal/app
}

func SetHostLanguage(lang string) {
	if _, ok := locales[lang]; ok {
		activeHostLang = lang
	}
}

func SetTargetLanguage(lang string) {
	activeTargetLang = lang
}

func GetHostLanguage() string {
	return activeHostLang
}

func GetTargetLanguage() string {
	return activeTargetLang
}

func GetLanguageName(lang string) string {
	names := map[string]string{
		"en": "English",
		"fr": "Français",
		"cs": "Čeština",
		"sk": "Slovenčina",
	}
	if name, ok := names[lang]; ok {
		return name
	}
	return lang
}

func T(id StringID, args ...any) string {
	translations, ok := locales[activeHostLang]
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
