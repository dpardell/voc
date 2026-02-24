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
	Short: "Download and install the French dictionary",
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		url := dictionary.KaikkiURL
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

		importer, err := dictionary.NewImporter(i18n.GetLanguage())
		if err != nil {
			fmt.Printf("Err: %v\n", err)
			return
		}
		if err := importer.DownloadAndImport(cmd.Context(), force, url); err != nil {
			fmt.Printf("Err: %v\n", err)
		}
	},
}
