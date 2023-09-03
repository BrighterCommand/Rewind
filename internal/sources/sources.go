package sources

import (
	"github.com/brightercommand/Rewind/internal/pages"
	"os"
)

const gitBookFileName = ".gitbook.yaml"
const tocFileName = ".toc.yaml"
const shardFolderName = "shared"
const summaryFolderName = "summary"
const sharedVersion = "Shared"

type Sources struct {
	Root     *pages.Root
	Shared   *pages.Shared
	Versions map[string]pages.Version
}

func NewSources() *Sources {
	return &Sources{
		Root:     &pages.Root{},
		Shared:   &pages.Shared{},
		Versions: make(map[string]pages.Version),
	}
}

// FindFromPath finds the sources for a book.
// It takes a root directory as a string.
// It returns a Sources struct.
// The Sources struct contains the root directory, shared documents, and versioned documents.
func (s *Sources) FindFromPath(root string) error {

	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	s.Root.SourcePath = root

	for _, entry := range entries {
		if entry.Name() == gitBookFileName {
			s.Root.GitBook = &pages.Doc{
				SourcePath: root,
				Storage:    entry,
			}
		} else if entry.IsDir() && entry.Name() == shardFolderName {
			shared, err := findShared(root, entry)
			if err != nil {
				return err
			}
			s.Shared = shared
		} else if entry.IsDir() && entry.Name() != summaryFolderName {
			version, err := findVersion(root, entry)
			if err != nil {
				return err
			}
			s.Versions[entry.Name()] = *version
		}
	}

	return err
}

// findVersionedDocs finds the versioned documents for the book.
// It takes a directory entry and an Version struct.
// It returns an error.
// It will recurse over any subdirectories and place all shared assets at the same level.
func findVersionedDocs(path string, version *pages.Version) (err error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if entry.Name() == tocFileName {
				version.TOC = &pages.Doc{SourcePath: path, Version: version.Version, Storage: entry}
			} else {
				version.Docs[entry.Name()] = pages.Doc{SourcePath: path, Version: version.Version, Storage: entry}
			}
		} else {
			err = findVersionedDocs(path+"/"+entry.Name(), version)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// findSharedDocs finds the documents for the book.
// It takes a directory entry and a Shared struct.
// It returns an error.
// It will recurse over any subdirectories and place all shared assets at the same level.
func findSharedDocs(path string, shared *pages.Shared) (err error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if entry.Name() == tocFileName {
				shared.TOC = &pages.Doc{SourcePath: path, Version: sharedVersion, Storage: entry}
			} else {
				shared.Docs[entry.Name()] = pages.Doc{SourcePath: path, Version: sharedVersion, Storage: entry}
			}
		} else {
			err = findSharedDocs(path+"/"+entry.Name(), shared)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// findShared finds the shared documents for the book.
// It takes a directory entry and its path.
// It returns a Shared struct.
func findShared(path string, entry os.DirEntry) (shared *pages.Shared, err error) {
	shared = &pages.Shared{
		Docs: make(map[string]pages.Doc),
	}

	sharedPath := path + "/" + entry.Name()

	err = findSharedDocs(sharedPath, shared)
	if err != nil {
		return shared, err
	}

	return shared, nil

}

// findVersion finds the versioned documents for the book.
// It takes a directory entry and its path.
// It returns a Version struct.
func findVersion(path string, entry os.DirEntry) (version *pages.Version, err error) {

	version = &pages.Version{
		Docs:    make(map[string]pages.Doc),
		Version: entry.Name(),
	}

	versionPath := path + "/" + entry.Name()

	err = findVersionedDocs(versionPath, version)
	if err != nil {
		return version, err
	}

	return version, nil

}