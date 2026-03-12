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
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		url := dictionary.GetDefaultKaikkiURL(vocApp.TargetLang)
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
				return fmt.Errorf("%s", i18n.T(i18n.InstallCancelled))
			}
			if input == "c" {
				fmt.Print(i18n.T(i18n.PromptCustomURL))
				input, _ = reader.ReadString('\n')
				url = strings.TrimSpace(input)
				if url == "" {
					return fmt.Errorf("%s", i18n.T(i18n.InstallCancelled))
				}
				break
			}
		}

		importer, err := dictionary.NewImporter(vocApp.TargetLang)
		if err != nil {
			return err
		}
		if err := importer.DownloadAndImport(cmd.Context(), force, url); err != nil {
			return err
		}
		return nil
	},
}
