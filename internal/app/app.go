package app

import (
	"fmt"
	"os"
	"voc/internal/config"
	"voc/internal/database"
	"voc/internal/dictionary"
	"voc/internal/i18n"
	"voc/internal/llm"
)

type App struct {
	DB         *database.Database
	Dict       *dictionary.Dictionary
	HostLang   string
	TargetLang string
	Settings   *config.Settings
}

func NewApp(hostLang, targetLang string) (*App, error) {
	settings, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: failed to load settings: %v\n", err)
	}

	if hostLang != "" {
		settings.HostLang = hostLang
	}
	if targetLang != "" {
		settings.TargetLang = targetLang
	}

	i18n.SetHostLanguage(settings.HostLang)
	i18n.SetTargetLanguage(settings.TargetLang)

	db, err := database.New()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", i18n.T(i18n.ErrDatabase), err)
	}

	dict, _ := dictionary.New(settings.TargetLang)

	return &App{
		DB:         db,
		Dict:       dict,
		HostLang:   settings.HostLang,
		TargetLang: settings.TargetLang,
		Settings:   settings,
	}, nil
}

func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}

func (a *App) ReinitDict() error {
	dict, err := dictionary.New(a.TargetLang)
	if err != nil {
		return err
	}
	a.Dict = dict
	return nil
}

func (a *App) GetLLMClient() (llm.LLMClient, error) {
	if mistralKey := os.Getenv("MISTRAL_API_KEY"); mistralKey != "" {
		return llm.NewMistralClient(mistralKey, os.Getenv("MISTRAL_MODEL")), nil
	}

	apiKey := os.Getenv("VERTEX_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("GEMINI_API_KEY")
	}

	projectID := os.Getenv("VERTEX_PROJECT_ID")
	location := os.Getenv("VERTEX_LOCATION")
	if location == "" {
		location = "us-central1"
	}

	modelName := os.Getenv("VERTEX_MODEL")

	if apiKey == "" || projectID == "" {
		return nil, fmt.Errorf("%s", i18n.T(i18n.LLMVarsMissing))
	}

	return llm.NewGeminiClient(apiKey, projectID, location, modelName)
}
