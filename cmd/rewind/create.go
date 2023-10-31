package rewind

import (
	"github.com/brightercommand/Rewind/internal/book"
	"github.com/brightercommand/Rewind/internal/sources"
	"github.com/spf13/cobra"
	"log"
)

var makeBookCmd = &cobra.Command{
	Use:     "makebook",
	Aliases: []string{"make"},
	Short:   "Turns versioned markdown files into a gitbook",
	Long: `Turns versioned markdown files into a gitbook.
			Expects a source path with a pointer to a director with the following structure:
			- source
                - shared //common docs
				    - .toc.yaml // table of contents
                    - doc.md // one or more markdown files
				- 1	// version number
				    - .toc.yaml // table of contents
                    - doc.md // one or more markdown files
                - 2 //version number
				    - .toc.yaml // table of contents
                    - doc.md // one or more markdown files
			Expects a destination path to write the gitbook directory to.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		log.Print("Creating book...")
		log.Print("Finding sources...")
		sources, err := findSources(args[0])
		if err != nil {
			log.Fatal(err)
		}

		log.Print("Making book...")
		book, err := book.MakeBook(sources, args[1])
		if err != nil {
			log.Fatal(err)
		}

		log.Print("Publishing book...")
		err = book.Publish()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func findSources(sourcePath string) (*sources.Sources, error) {
	log.Print("Finding sources in " + sourcePath + "...")

	src := sources.NewSources()
	err := src.FindFromPath(sourcePath)
	if err != nil {
		return nil, err
	}
	return src, nil
}
