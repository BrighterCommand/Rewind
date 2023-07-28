package internal

import (
	"io"
	"os"
)

type Shared struct {
	Docs map[string]Doc
	TOC  Doc
}

type Version struct {
	DestPath string
	Docs     map[string]Doc
	TOC      Doc
	Version  string
}

type Doc struct {
	SourcePath string
	Version    string
	Storage    os.DirEntry
}

type Root struct {
	DestPath string
	Summary  Doc
	GitBook  Doc
}

type Book struct {
	Root     *Root
	Versions map[string]Version
}

func (b *Book) Create() error {

	//create root directory
	rootPath := b.Root.DestPath
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		err := os.MkdirAll(rootPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//copy root files
	err := copyFile(b.Root.Summary.SourcePath, b.Root.DestPath, b.Root.Summary.Storage.Name())
	if err != nil {
		return err
	}

	err = copyFile(b.Root.GitBook.SourcePath, b.Root.DestPath, b.Root.GitBook.Storage.Name())
	if err != nil {
		return err
	}

	//copy versioned files
	for _, version := range b.Versions {

		destPath := version.DestPath
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			err := os.MkdirAll(destPath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		for _, doc := range version.Docs {
			err = copyFile(doc.SourcePath, destPath, doc.Storage.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Sources struct {
	Root     *Root
	Shared   *Shared
	Versions map[string]Version
}

func (s *Sources) BuildBook(destPath string) *Book {
	book := &Book{
		Root: &Root{
			DestPath: destPath,
			Summary:  s.Root.Summary,
			GitBook:  s.Root.GitBook,
		},
		Versions: make(map[string]Version),
	}

	for key, version := range s.Versions {

		var bookVersion = &Version{
			Version:  version.Version,
			DestPath: destPath + "/" + version.Version,
			Docs:     make(map[string]Doc),
		}

		//copy shared assets first
		for key, doc := range s.Shared.Docs {
			bookVersion.Docs[key] = doc
		}

		//now copy assets for this version and overwrite any shared assets with the same name
		for key, doc := range version.Docs {
			bookVersion.Docs[key] = doc
		}

		book.Versions[key] = *bookVersion
	}

	return book
}

// FindSources finds the sources for a book.
// It takes a root directory as a string.
// It returns a Sources struct.
// The Sources struct contains the root directory, shared documents, and versioned documents.
func FindSources(root string) (sources *Sources, err error) {

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	sources = &Sources{
		Root:     &Root{},
		Shared:   &Shared{},
		Versions: make(map[string]Version),
	}

	for _, entry := range entries {
		if entry.Name() == "SUMMARY.md" {
			sources.Root.Summary = Doc{
				SourcePath: root,
				Storage:    entry,
			}
		} else if entry.Name() == ".gitbook.yaml" {
			sources.Root.GitBook = Doc{
				SourcePath: root,
				Storage:    entry,
			}
		} else if entry.IsDir() && entry.Name() == "shared" {
			shared, err := findShared(root, entry)
			if err != nil {
				return nil, err
			}
			sources.Shared = shared
		} else if entry.IsDir() {
			version, err := findVersion(root, entry)
			if err != nil {
				return nil, err
			}
			sources.Versions[entry.Name()] = *version
		}
	}

	return sources, err
}

// copyFile copies a file from sourcePath to destPath
// It takes two string arguments: sourcePath and destPath
// It creates the destination file if it does not exist
// It overwrites the destination file if it exists
// It returns an error if the copy fails
func copyFile(sourcePath string, destPath string, fileName string) (err error) {
	r, err := os.Open(sourcePath + "/" + fileName)
	if err != nil {
		return err
	}
	defer r.Close()

	//if the directory does not exist, create it
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		err = os.MkdirAll(destPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	w, err := os.Create(destPath + "/" + fileName)
	if err != nil {
		return err
	}

	defer func() {
		if c := w.Close(); c != nil && err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)
	return err
}

// findVersionedDocs finds the versioned documents for the book.
// It takes a directory entry and an Version struct.
// It returns an error.
// It will recurse over any subdirectories and place all shared assets at the same level.
func findVersionedDocs(path string, version *Version) (err error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			version.Docs[entry.Name()] = Doc{SourcePath: path, Version: version.Version, Storage: entry}
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
func findSharedDocs(path string, shared *Shared) (err error) {

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			shared.Docs[entry.Name()] = Doc{SourcePath: path, Version: "Shared", Storage: entry}
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
func findShared(path string, entry os.DirEntry) (shared *Shared, err error) {
	shared = &Shared{
		Docs: make(map[string]Doc),
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
func findVersion(path string, entry os.DirEntry) (version *Version, err error) {

	version = &Version{
		Docs:    make(map[string]Doc),
		Version: entry.Name(),
	}

	versionPath := path + "/" + entry.Name()

	err = findVersionedDocs(versionPath, version)
	if err != nil {
		return version, err
	}

	return version, nil

}
