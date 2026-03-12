package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"voc/internal/app"
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
	rootCmd.PersistentFlags().StringVarP(&hostLangFlag, "host-lang", "l", "", "UI language (en, fr)")
	rootCmd.PersistentFlags().StringVarP(&targetLangFlag, "target-lang", "t", "", "Language you are learning (fr, es, etc.)")
}

func runSplashScreen(ctx context.Context) {
	var lastErr string

	for {
		option, err := ui.RunSplash(lastErr, vocApp.Dict != nil)
		lastErr = ""
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		switch option {
		case ui.OptionSearch:
			searchCmd.SetContext(ctx)
			if err := searchCmd.RunE(searchCmd, nil); err != nil {
				lastErr = err.Error()
			}
		case ui.OptionQuiz:
			quizCmd.SetContext(ctx)
			if err := quizCmd.RunE(quizCmd, nil); err != nil {
				lastErr = err.Error()
			}
		case ui.OptionConvo:
			convoCmd.SetContext(ctx)
			if err := convoCmd.RunE(convoCmd, nil); err != nil {
				lastErr = err.Error()
			}
		case ui.OptionList:
			listCmd.SetContext(ctx)
			if err := listCmd.RunE(listCmd, nil); err != nil {
				lastErr = err.Error()
			}
		case ui.OptionQuit:
			return
		default:
			return
		}
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
