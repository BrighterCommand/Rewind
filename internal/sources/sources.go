package sources

import (
	"github.com/brightercommand/Rewind/internal/pages"
	"log"
	"os"
	"strings"
)

const gitBookFileName = ".gitbook.yaml"
const tocFileName = ".toc.yaml"
const sharedFolderName = "shared"
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
		} else if entry.Name() == "README.md" {
			s.Root.ReadMe = &pages.Doc{
				SourcePath: root,
				Storage:    entry,
			}
		} else if entry.IsDir() && entry.Name() == sharedFolderName {
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
// It takes a directory entry and a Version struct.
// It returns an error.
// We assume that all documents are in the root of the version folder.
// We assume that sub-folders are used for images.
func findVersionedDocs(path string, version *pages.Version) (err error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if entry.Name() == tocFileName {
				version.TOC = &pages.Doc{SourcePath: path, Version: version.Version, Storage: entry}
			} else if isMarkDownFile(entry) {
				version.Docs[entry.Name()] = pages.Doc{SourcePath: path, Version: version.Version, Storage: entry}
			}
		} else if entry.Name() == pages.StaticFolderName {
			err = findVersionedImages(path+"/"+entry.Name(), version)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func findVersionedImages(path string, version *pages.Version) error {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if isImageFile(entry.Name()) {
				version.Images[entry.Name()] = pages.Asset{SourcePath: path, What: pages.Image, Version: version.Version, Storage: entry}
			}
		} else if entry.Name() == pages.ImageFolderName {
			err = findVersionedImages(path+"/"+entry.Name(), version)
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
// We assume that the documents are the root level
// We assume that sub-folders are used for images.
func findSharedDocs(path string, shared *pages.Shared) (err error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if entry.Name() == tocFileName {
				shared.TOC = &pages.Doc{SourcePath: path, Version: sharedVersion, Storage: entry}
			} else if isMarkDownFile(entry) {
				shared.Docs[entry.Name()] = pages.Doc{SourcePath: path, Version: sharedVersion, Storage: entry}
			}
		} else if entry.Name() == pages.StaticFolderName {
			err = findSharedImages(path+"/"+entry.Name(), shared)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func findSharedImages(path string, shared *pages.Shared) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if isImageFile(entry.Name()) {
				shared.Images[entry.Name()] = pages.Asset{SourcePath: path, What: pages.Image, Version: sharedVersion, Storage: entry}
			}
		} else if entry.Name() == pages.ImageFolderName {
			err = findSharedImages(path+"/"+entry.Name(), shared)
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
		Docs:   make(map[string]pages.Doc),
		Images: make(map[string]pages.Asset),
	}

	sharedPath := path + "/" + entry.Name()

	log.Print("Finding shared docs in " + sharedPath + "...")
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
		Images:  make(map[string]pages.Asset),
		Version: entry.Name(),
	}

	versionPath := path + "/" + entry.Name()
	log.Print("Finding versioned docs in " + versionPath + "...")

	err = findVersionedDocs(versionPath, version)
	if err != nil {
		return version, err
	}

	return version, nil

}

// Helper function to check if a file has an image extension
func isImageFile(filename string) bool {
	// Define a list of image file extensions
	imageExtensions := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".svg"}

	// Check if the filename ends with any of the image extensions
	for _, ext := range imageExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func isMarkDownFile(entry os.DirEntry) bool {
	return strings.HasSuffix(entry.Name(), ".md")
}
