package rewind

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "rewind",
	Short: "rewind - a simple CLI to create and publish Brighter's GitBook documentation",
	Long: `rewind - takes both shared and version specific documentation and publishes it to a single GitBook
			repository. It also creates a summary.md file that is used by GitBook to create the table of contents.
			Note that GitBook integration with the repo is assumed, via a configured root folder. This root
			folder should be your destination folder. The source folder should also be a root level folder
			which contains a shared folder for common markdown pages, and version folders for versioned documents.
			A toc file in shared and version is used to supply information on how to build the toc.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(makeBookCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
