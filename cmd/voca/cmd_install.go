package main

import (
	"fmt"
	"voca/internal/dictionary"

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
		importer, err := dictionary.NewImporter(detectedLang)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		if err := importer.DownloadAndImport(force); err != nil {
			fmt.Printf("Error installing dictionary: %v\n", err)
		}
	},
}
