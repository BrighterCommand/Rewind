package rewind

import (
	"github.com/brightercommand/Rewind/internal/book"
	"github.com/brightercommand/Rewind/internal/sources"
	"github.com/spf13/cobra"
	"log"
)

var reverseCmd = &cobra.Command{
	Use:     "makebook",
	Aliases: []string{"make"},
	Short:   "Turns versioned markdown files into a gitbook director",
	Long: `Turns versioned markdown files into a gitbook directory.
			Expects a source path with a pointer to a director with the following structure:
			- source
                - shared //common docs
				    - .toc.yaml // table of contents
                    - doc.md // one or more markdown files
				- 1.0.0	// version number
				    - .toc.yaml // table of contents
                    - doc.md // one or more markdown files
                - 2.0.0 //version number
				    - .toc.yaml // table of contents
                    - doc.md // one or more markdown files
			Expects a destination path to write the gitbook directory to.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

		sources, err := findSources(args[0])
		if err != nil {
			log.Fatal(err)
		}

		book, err := book.MakeBook(sources, args[1])
		if err != nil {
			log.Fatal(err)
		}

		err = book.Publish()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func findSources(sourcePath string) (*sources.Sources, error) {
	src := sources.NewSources()
	err := src.FindFromPath(sourcePath)
	if err != nil {
		return nil, err
	}
	return src, nil
}
