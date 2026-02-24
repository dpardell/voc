package app

import (
	"fmt"
	"os"
	"voc/internal/database"
	"voc/internal/dictionary"
	"voc/internal/i18n"
	"voc/internal/llm"
)

type App struct {
	DB   *database.Database
	Dict *dictionary.Dictionary
	Lang string
}

func NewApp(lang string) (*App, error) {
	db, err := database.New()
	if err != nil {
		return nil, fmt.Errorf("%s: %v", i18n.T(i18n.ErrDatabase), err)
	}

	dict, _ := dictionary.New(lang)

	return &App{
		DB:   db,
		Dict: dict,
		Lang: lang,
	}, nil
}

func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}

func (a *App) GetLLMClient() (*llm.Client, error) {
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
		return nil, fmt.Errorf("VERTEX_API_KEY and VERTEX_PROJECT_ID must be set")
	}

	return llm.NewClient(apiKey, projectID, location, modelName)
}
