package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"voc/internal/dictionary"
	"voc/internal/i18n"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(installDictCmd)
	installDictCmd.Flags().Bool("force", false, "Force re-download even if cached")
}

var installDictCmd = &cobra.Command{
	Use:   "install-dict",
	Short: "Download and install the dictionary for the target language",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		url := dictionary.KaikkiURL
		if vocApp.Settings != nil && vocApp.Settings.Dictionaries != nil {
			if dURL, ok := vocApp.Settings.Dictionaries[vocApp.TargetLang]; ok {
				url = dURL
			}
		}

		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Printf(i18n.T(i18n.PromptDownloadDefault), url)
			input, _ := reader.ReadString('\n')
			input = strings.ToLower(strings.TrimSpace(input))

			if input == "" || input == "y" {
				break
			}
			if input == "n" {
				fmt.Println(i18n.T(i18n.InstallCancelled))
				return
			}
			if input == "c" {
				fmt.Print(i18n.T(i18n.PromptCustomURL))
				input, _ = reader.ReadString('\n')
				url = strings.TrimSpace(input)
				if url == "" {
					fmt.Println(i18n.T(i18n.InstallCancelled))
					return
				}
				break
			}
		}

		importer, err := dictionary.NewImporter(vocApp.TargetLang)
		if err != nil {
			fmt.Printf("Err: %v\n", err)
			return
		}
		if err := importer.DownloadAndImport(cmd.Context(), force, url); err != nil {
			fmt.Printf("Err: %v\n", err)
		}
	},
}
