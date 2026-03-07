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
	NoWordsToQuiz  StringID = "NoWordsToQuiz"
	LLMVarsMissing StringID = "LLMVarsMissing"
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
	MenuVocab:              "MY VOCABULARY",
	MenuInstall:            "INSTALL DICTIONARY",
	MenuQuit:               "QUIT",
	SplashFooter:           "↑/↓: navigate • enter: select • q: quit",
	ImportComplete:         "Import complete: %d added, %d skipped",
	ImportUsage:            "Usage: voc import -f <file> OR pipe input",
	ExportSuccess:          "Exported %d words to %s",
	QuizGenerating:         "Generating quiz...",
	QuizUpdating:           "Updating progress...",
	QuizHintType:           "Type your answer and press Enter",
	NoWordsToQuiz:          "No words to quiz.",
	LLMVarsMissing:         "Set MISTRAL_API_KEY, or set VERTEX_API_KEY (or GEMINI_API_KEY) and VERTEX_PROJECT_ID",
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
	MenuVocab:              "MON VOCABULAIRE",
	MenuInstall:            "INSTALLER LE DICTIONNAIRE",
	MenuQuit:               "QUITTER",
	SplashFooter:           "↑/↓: naviguer • entrée: sélectionner • q: quitter",
	ImportComplete:         "Importation terminée : %d ajouté(s), %d ignoré(s)",
	ImportUsage:            "Utilisation : voc import -f <fichier> OU pipe entrée",
	ExportSuccess:          "Exportation de %d mots vers %s",
	QuizGenerating:         "Génération du quiz...",
	QuizUpdating:           "Mise à jour de la progression...",
	QuizHintType:           "Tapez votre réponse et appuyez sur Entrée",
	NoWordsToQuiz:          "Aucun mot pour le quiz.",
	LLMVarsMissing:         "Définissez MISTRAL_API_KEY, ou VERTEX_API_KEY (ou GEMINI_API_KEY) et VERTEX_PROJECT_ID",
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
	NoWordsToQuiz:          "Žádná slova pro kvíz.",
	LLMVarsMissing:         "Nastavte MISTRAL_API_KEY, nebo VERTEX_API_KEY (nebo GEMINI_API_KEY) a VERTEX_PROJECT_ID",
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
	NoWordsToQuiz:          "Žiadne slová pre kvíz.",
	LLMVarsMissing:         "Nastavte MISTRAL_API_KEY, alebo VERTEX_API_KEY (alebo GEMINI_API_KEY) a VERTEX_PROJECT_ID",
}

var es = map[StringID]string{
	NoMatches:              "No se encontraron coincidencias...",
	SelectItem:             "Seleccione un elemento para ver los detalles...",
	ErrorPreview:           "Error al cargar la vista previa",
	ResizeWindow:           "(Cambie el tamaño de la ventana para ver los resultados)",
	TypeToSearch:           "Escriba para buscar...",
	Searching:              "BUSCAR EN EL DICCIONARIO",
	AddingWord:             "AÑADIR PALABRA",
	FixingWord:             "CORRIGIENDO PALABRA: %s",
	ShowingWord:            "MOSTRAR PALABRA",
	DictionaryNotInstalled: "Diccionario no instalado. Ejecute 'voc install-dict' primero.",
	ErrWordNotFound:        "Palabra no encontrada en el diccionario",
	ErrDatabase:            "Error al inicializar la base de datos",
	DictionaryNotFound:     "Diccionario no encontrado en %s",
	PromptDownloadDefault:  "¿Descargar diccionario de %s? [Y/n/c (personalizado)] ",
	PromptCustomURL:        "Ingrese la URL del diccionario personalizado: ",
	InstallCancelled:       "Instalación cancelada.",
	HintsSearch:            "Ctrl-C: Salir • Ctrl-P: Arriba • Enter: Seleccionar",
	HintsDefinition:        "Esc: Volver • J/K: Desplazar • Ctrl-P/N: Ant/Sig • Ctrl-S: Alternar Vocabulario",
	Saved:                  "GUARDADO ❤️",
	NotSaved:               "NO GUARDADO",
	WordCount:              "Tu lista (%d palabras)",
	QuizTitle:              "CUESTIONARIO DE VOCABULARIO",
	QuizQuestion:           "Pregunta %d/%d",
	QuizScore:              "Puntuación: %d/%d",
	QuizCorrect:            "¡CORRECTO! ✨",
	QuizWrong:              "¡INCORRECTO! La respuesta correcta era: %s",
	QuizCompleted:          "¡Cuestionario completado! Puntuación final: %d/%d",
	HintsQuiz:              "Enter: Enviar/Siguiente • Ctrl-C: Salir",
	ConvoPlaceholder:       "Diga algo en su idioma de destino...",
	ConvoStarted:           "¡Chat iniciado! ¡Diga hola!",
	ConvoCorrections:       "Correcciones:",
	ConvoThinking:          "Pensando...",
	ConvoError:             "ERROR",
	ConvoExitHint:          "Presione Ctrl+C para salir",
	ConvoCoachTitle:        "ENTRENADOR DE IDIOMAS",
	QuizPlaceholder:        "Escriba su respuesta...",
	PressEnterToCont:       "Presione Enter para continuar...",
	ErrorNoWords:           "No hay suficientes palabras para generar un cuestionario. ¡Añada más palabras primero!",
	ErrorNoSavedWords:      "Aún no ha guardado ninguna palabra. ¡Busque palabras para añadirlas!",
	DictionaryLoading:      "Diccionario no cargado",
	SplashTagline:          "Su compañero de idiomas con IA",
	MenuSearch:             "BUSCAR EN EL DICCIONARIO",
	MenuQuiz:               "CUESTIONARIO DE VOCABULARIO",
	MenuConvo:              "ENTRENADOR DE IDIOMAS",
	MenuVocab:              "MI VOCABULARIO",
	MenuInstall:            "INSTALAR DICCIONARIO",
	MenuQuit:               "SALIR",
	SplashFooter:           "↑/↓: navegar • enter: seleccionar • q: salir",
	ImportComplete:         "Importación completada: %d añadidas, %d omitidas",
	ImportUsage:            "Uso: voc import -f <archivo> O entrada por tubería",
	ExportSuccess:          "Exportadas %d palabras a %s",
	QuizGenerating:         "Generando cuestionario...",
	QuizUpdating:           "Actualizando progreso...",
	QuizHintType:           "Escriba su respuesta y presione Enter",
	NoWordsToQuiz:          "No hay palabras para el cuestionario.",
	LLMVarsMissing:         "Configure MISTRAL_API_KEY, o configure VERTEX_API_KEY (o GEMINI_API_KEY) y VERTEX_PROJECT_ID",
}

var de = map[StringID]string{
	NoMatches:              "Keine Treffer gefunden...",
	SelectItem:             "Wählen Sie ein Element aus, um Details zu sehen...",
	ErrorPreview:           "Fehler beim Laden der Vorschau",
	ResizeWindow:           "(Fenstergröße ändern, um Ergebnisse zu sehen)",
	TypeToSearch:           "Tippen Sie, um zu suchen...",
	Searching:              "WÖRTERBUCH DURCHSUCHEN",
	AddingWord:             "WORT HINZUFÜGEN",
	FixingWord:             "WORT KORRIGIEREN: %s",
	ShowingWord:            "WORT ANZEIGEN",
	DictionaryNotInstalled: "Wörterbuch nicht installiert. Führen Sie zuerst 'voc install-dict' aus.",
	ErrWordNotFound:        "Wort nicht im Wörterbuch gefunden",
	ErrDatabase:            "Fehler beim Initialisieren der Datenbank",
	DictionaryNotFound:     "Wörterbuch nicht gefunden unter %s",
	PromptDownloadDefault:  "Wörterbuch von %s herunterladen? [Y/n/c (benutzerdefiniert)] ",
	PromptCustomURL:        "Geben Sie eine benutzerdefinierte Wörterbuch-URL ein: ",
	InstallCancelled:       "Installation abgebrochen.",
	HintsSearch:            "Strg-C: Beenden • Strg-P: Hoch • Enter: Auswählen",
	HintsDefinition:        "Esc: Zurück • J/K: Scrollen • Strg-P/N: Zurück/Vor • Strg-S: Wortschatz umschalten",
	Saved:                  "GESPEICHERT ❤️",
	NotSaved:               "NICHT GESPEICHERT",
	WordCount:              "Ihre Liste (%d Wörter)",
	QuizTitle:              "VOKABELQUIZ",
	QuizQuestion:           "Frage %d/%d",
	QuizScore:              "Ergebnis: %d/%d",
	QuizCorrect:            "RICHTIG! ✨",
	QuizWrong:              "FALSCH! Die richtige Antwort war: %s",
	QuizCompleted:          "Quiz abgeschlossen! Endergebnis: %d/%d",
	HintsQuiz:              "Enter: Absenden/Weiter • Strg-C: Beenden",
	ConvoPlaceholder:       "Sagen Sie etwas in Ihrer Zielsprache...",
	ConvoStarted:           "Chat gestartet! Sagen Sie Hallo!",
	ConvoCorrections:       "Korrekturen:",
	ConvoThinking:          "Nachdenken...",
	ConvoError:             "FEHLER",
	ConvoExitHint:          "Drücken Sie Strg+C zum Beenden",
	ConvoCoachTitle:        "SPRACHCOACH",
	QuizPlaceholder:        "Geben Sie Ihre Antwort ein...",
	PressEnterToCont:       "Drücken Sie Enter, um fortzufahren...",
	ErrorNoWords:           "Nicht genügend Wörter, um ein Quiz zu erstellen. Fügen Sie zuerst weitere Wörter hinzu!",
	ErrorNoSavedWords:      "Sie haben noch keine Wörter gespeichert. Suchen Sie nach Wörtern, um sie hinzuzufügen!",
	DictionaryLoading:      "Wörterbuch nicht geladen",
	SplashTagline:          "Ihr KI-gestützter Sprachbegleiter",
	MenuSearch:             "WÖRTERBUCH DURCHSUCHEN",
	MenuQuiz:               "VOKABELQUIZ",
	MenuConvo:              "SPRACHCOACH",
	MenuVocab:              "MEIN WORTSCHATZ",
	MenuInstall:            "WÖRTERBUCH INSTALLIEREN",
	MenuQuit:               "BEENDEN",
	SplashFooter:           "↑/↓: Navigieren • Enter: Auswählen • q: Beenden",
	ImportComplete:         "Import abgeschlossen: %d hinzugefügt, %d übersprungen",
	ImportUsage:            "Verwendung: voc import -f <Datei> ODER Pipe-Eingabe",
	ExportSuccess:          "%d Wörter nach %s exportiert",
	QuizGenerating:         "Quiz wird erstellt...",
	QuizUpdating:           "Fortschritt wird aktualisiert...",
	QuizHintType:           "Geben Sie Ihre Antwort ein und drücken Sie Enter",
	NoWordsToQuiz:          "Keine Wörter für das Quiz.",
	LLMVarsMissing:         "Setzen Sie MISTRAL_API_KEY, oder VERTEX_API_KEY (oder GEMINI_API_KEY) und VERTEX_PROJECT_ID",
}

var pt = map[StringID]string{
	NoMatches:              "Nenhuma correspondência encontrada...",
	SelectItem:             "Selecione um item para ver os detalhes...",
	ErrorPreview:           "Erro ao carregar a pré-visualização",
	ResizeWindow:           "(Redimensione a janela para ver os resultados)",
	TypeToSearch:           "Digite para pesquisar...",
	Searching:              "PESQUISAR NO DICIONÁRIO",
	AddingWord:             "ADICIONAR PALAVRA",
	FixingWord:             "A CORRIGIR PALAVRA: %s",
	ShowingWord:            "MOSTRAR PALABRA",
	DictionaryNotInstalled: "Dicionário não instalado. Execute 'voc install-dict' primeiro.",
	ErrWordNotFound:        "Palavra não encontrada no dicionário",
	ErrDatabase:            "Erro ao inicializar a base de dados",
	DictionaryNotFound:     "Dicionário não encontrado em %s",
	PromptDownloadDefault:  "Descarregar dicionário de %s? [Y/n/c (personalizado)] ",
	PromptCustomURL:        "Introduza o URL do dicionário personalizado: ",
	InstallCancelled:       "Instalação cancelada.",
	HintsSearch:            "Ctrl-C: Sair • Ctrl-P: Cima • Enter: Selecionar",
	HintsDefinition:        "Esc: Voltar • J/K: Deslocar • Ctrl-P/N: Ant/Seg • Ctrl-S: Alternar Vocabulário",
	Saved:                  "GUARDADO ❤️",
	NotSaved:               "NÃO GUARDADO",
	WordCount:              "A sua lista (%d palavras)",
	QuizTitle:              "QUESTIONÁRIO DE VOCABULÁRIO",
	QuizQuestion:           "Pergunta %d/%d",
	QuizScore:              "Pontuação: %d/%d",
	QuizCorrect:            "CORRETO! ✨",
	QuizWrong:              "ERRADO! A resposta correta era: %s",
	QuizCompleted:          "Questionário concluído! Pontuação final: %d/%d",
	HintsQuiz:              "Enter: Enviar/Seguinte • Ctrl-C: Sair",
	ConvoPlaceholder:       "Diga algo na sua língua de destino...",
	ConvoStarted:           "Chat iniciado! Diga olá!",
	ConvoCorrections:       "Correções:",
	ConvoThinking:          "A pensar...",
	ConvoError:             "ERRO",
	ConvoExitHint:          "Pressione Ctrl+C para sair",
	ConvoCoachTitle:        "TREINADOR DE LÍNGUAS",
	QuizPlaceholder:        "Escreva a sua resposta...",
	PressEnterToCont:       "Pressione Enter para continuar...",
	ErrorNoWords:           "Não há palavras suficientes para gerar um questionário. Adicione mais palavras primeiro!",
	ErrorNoSavedWords:      "Ainda não guardou nenhuma palavra. Pesquise palavras para as adicionar!",
	DictionaryLoading:      "Dicionário não carregado",
	SplashTagline:          "O seu companheiro de línguas com IA",
	MenuSearch:             "PESQUISAR NO DICIONÁRIO",
	MenuQuiz:               "QUESTIONÁRIO DE VOCABULÁRIO",
	MenuConvo:              "TREINADOR DE LÍNGUAS",
	MenuVocab:              "O MEU VOCABULÁRIO",
	MenuInstall:            "INSTALAR DICIONÁRIO",
	MenuQuit:               "SAIR",
	SplashFooter:           "↑/↓: navegar • enter: selecionar • q: sair",
	ImportComplete:         "Importação concluída: %d adicionadas, %d ignoradas",
	ImportUsage:            "Uso: voc import -f <ficheiro> OU entrada por canal",
	ExportSuccess:          "Exportadas %d palavras para %s",
	QuizGenerating:         "A gerar questionário...",
	QuizUpdating:           "A atualizar progresso...",
	QuizHintType:           "Escreva a sua resposta e pressione Enter",
	NoWordsToQuiz:          "Sem palavras para o questionário.",
	LLMVarsMissing:         "Defina MISTRAL_API_KEY, ou VERTEX_API_KEY (ou GEMINI_API_KEY) e VERTEX_PROJECT_ID",
}

var ptBR = map[StringID]string{
	NoMatches:              "Nenhuma correspondência encontrada...",
	SelectItem:             "Selecione um item para ver detalhes...",
	ErrorPreview:           "Erro ao carregar a pré-visualização",
	ResizeWindow:           "(Redimensione a janela para ver os resultados)",
	TypeToSearch:           "Digite para pesquisar...",
	Searching:              "PESQUISAR NO DICIONÁRIO",
	AddingWord:             "ADICIONAR PALAVRA",
	FixingWord:             "CORRIGINDO PALAVRA: %s",
	ShowingWord:            "MOSTRAR PALAVRA",
	DictionaryNotInstalled: "Dicionário não instalado. Execute 'voc install-dict' primeiro.",
	ErrWordNotFound:        "Palavra não encontrada no dicionário",
	ErrDatabase:            "Erro ao inicializar o banco de dados",
	DictionaryNotFound:     "Dicionário não encontrado em %s",
	PromptDownloadDefault:  "Baixar dicionário de %s? [Y/n/c (personalizado)] ",
	PromptCustomURL:        "Insira a URL do dicionário personalizado: ",
	InstallCancelled:       "Instalação cancelada.",
	HintsSearch:            "Ctrl-C: Sair • Ctrl-P: Cima • Enter: Selecionar",
	HintsDefinition:        "Esc: Voltar • J/K: Rolar • Ctrl-P/N: Ant/Próx • Ctrl-S: Alternar Vocabulário",
	Saved:                  "SALVO ❤️",
	NotSaved:               "NÃO SALVO",
	WordCount:              "Sua lista (%d palavras)",
	QuizTitle:              "QUIZ DE VOCABULÁRIO",
	QuizQuestion:           "Pergunta %d/%d",
	QuizScore:              "Pontuação: %d/%d",
	QuizCorrect:            "CORRETO! ✨",
	QuizWrong:              "ERRADO! A resposta correta era: %s",
	QuizCompleted:          "Quiz concluído! Pontuação final: %d/%d",
	HintsQuiz:              "Enter: Enviar/Próxima • Ctrl-C: Sair",
	ConvoPlaceholder:       "Diga algo no seu idioma de destino...",
	ConvoStarted:           "Chat iniciado! Diga olá!",
	ConvoCorrections:       "Correções:",
	ConvoThinking:          "Pensando...",
	ConvoError:             "ERRO",
	ConvoExitHint:          "Pressione Ctrl+C para sair",
	ConvoCoachTitle:        "TREINADOR DE IDIOMAS",
	QuizPlaceholder:        "Escreva sua resposta...",
	PressEnterToCont:       "Pressione Enter para continuar...",
	ErrorNoWords:           "Não há palavras suficientes para gerar um quiz. Adicione mais palavras primeiro!",
	ErrorNoSavedWords:      "Você ainda não salvou nenhuma palavra. Pesquise palavras para adicioná-las!",
	DictionaryLoading:      "Dicionário não carregado",
	SplashTagline:          "Seu companheiro de idiomas com IA",
	MenuSearch:             "PESQUISAR NO DICIONÁRIO",
	MenuQuiz:               "QUIZ DE VOCABULÁRIO",
	MenuConvo:              "TREINADOR DE IDIOMAS",
	MenuVocab:              "MEU VOCABULÁRIO",
	MenuInstall:            "INSTALAR DICIONÁRIO",
	MenuQuit:               "SAIR",
	SplashFooter:           "↑/↓: navegar • enter: selecionar • q: sair",
	ImportComplete:         "Importação concluída: %d adicionadas, %d ignoradas",
	ImportUsage:            "Uso: voc import -f <arquivo> OU entrada por pipe",
	ExportSuccess:          "Exportadas %d palavras para %s",
	QuizGenerating:         "Gerando quiz...",
	QuizUpdating:           "Atualizando progresso...",
	QuizHintType:           "Escreva sua resposta e pressione Enter",
	NoWordsToQuiz:          "Sem palavras para o quiz.",
	LLMVarsMissing:         "Defina MISTRAL_API_KEY, ou VERTEX_API_KEY (ou GEMINI_API_KEY) e VERTEX_PROJECT_ID",
}

var activeHostLang = "en"
var activeTargetLang = "fr"

var locales = map[string]map[StringID]string{
	"en":    en,
	"fr":    fr,
	"cs":    cs,
	"sk":    sk,
	"es":    es,
	"de":    de,
	"pt":    pt,
	"pt-br": ptBR,
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
		"en":    "English",
		"fr":    "Français",
		"cs":    "Čeština",
		"sk":    "Slovenčina",
		"es":    "Español",
		"de":    "Deutsch",
		"pt":    "Português",
		"pt-br": "Português (Brasil)",
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
