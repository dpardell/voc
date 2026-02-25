package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"voc/internal/app"
	"voc/internal/database"
	"voc/internal/i18n"
	"voc/internal/ui"
)

var vocApp *app.App
var hostLangFlag string
var targetLangFlag string

var rootCmd = &cobra.Command{
	Use:   "voc",
	Short: "Voc - Vocabulary trainer",
	Long: `Voc tracks the words you learn, fetches definitions from Wiktionary,
and helps you review with interactive quizzes.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		vocApp, err = app.NewApp(hostLangFlag, targetLangFlag)
		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if vocApp != nil {
			return vocApp.Close()
		}
		return nil
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			runSplashScreen(cmd.Context())
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&hostLangFlag, "host-lang", "", "UI language (en, fr)")
	rootCmd.PersistentFlags().StringVar(&targetLangFlag, "target-lang", "", "Language you are learning (fr, es, etc.)")
}

type DailySaying struct {
	Date   string `json:"date"`
	Saying string `json:"saying"`
}

func getDailySaying(ctx context.Context) string {
	cachePath, _ := database.GetDailySayingPath()
	today := time.Now().Format("2006-01-02")

	// Try reading cache
	data, err := os.ReadFile(cachePath)
	if err == nil {
		var ds DailySaying
		if json.Unmarshal(data, &ds) == nil && ds.Date == today {
			return ds.Saying
		}
	}

	// Generate new saying
	client, err := vocApp.GetLLMClient()
	if err != nil {
		return i18n.T(i18n.DailySayingPrompt)
	}
	defer client.Close()

	progress, _ := database.GetProgress()

	saying, err := client.GenerateDailySaying(ctx, vocApp.TargetLang, progress)
	if err != nil {
		return "Time to learn some new words!"
	}

	// Save to cache
	ds := DailySaying{Date: today, Saying: saying}
	dsBytes, _ := json.Marshal(ds)
	os.WriteFile(cachePath, dsBytes, 0644)

	return saying
}

func runSplashScreen(ctx context.Context) {
	saying := getDailySaying(ctx)

	for {
		option, err := ui.RunSplash(saying)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch option {
		case ui.OptionSearch:
			searchCmd.SetContext(ctx)
			searchCmd.Run(searchCmd, nil)
		case ui.OptionQuiz:
			quizCmd.SetContext(ctx)
			quizCmd.Run(quizCmd, nil)
		case ui.OptionConvo:
			convoCmd.SetContext(ctx)
			convoCmd.Run(convoCmd, nil)
		case ui.OptionList:
			listCmd.SetContext(ctx)
			listCmd.Run(listCmd, nil)
		case ui.OptionQuit:
			return
		default:
			return
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
