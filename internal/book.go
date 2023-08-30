package internal

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
)

type Book struct {
	Root     *Root
	Versions map[string]Version
}

func newBook(destPath string, sourcePath string) *Book {
	return &Book{
		Root: &Root{
			DestPath:   destPath,
			SourcePath: sourcePath,
		},
		Versions: make(map[string]Version),
	}
}

func (b *Book) Publish() error {

	//create root directory
	rootPath := b.Root.DestPath
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		err := os.MkdirAll(rootPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	//copy root files
	err := copyFile(b.Root.WorkDir, b.Root.DestPath, SummaryFileName)
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

	//clean up the temporary files
	err = b.ClearWorkDir()
	if err != nil {
		return err
	}

	return nil
}

func (b *Book) MakeVersions(s *Sources) {
	for key, version := range s.Versions {

		var bookVersion = &Version{
			Version:  version.Version,
			DestPath: b.Root.DestPath + "/" + ContentDirName + "/" + version.Version,
			Docs:     make(map[string]Doc),
		}

		//copy shared assets first
		for key, doc := range s.Shared.Docs {
			bookVersion.Docs[key] = doc
		}

		//now copy assets for this sharedVersion and overwrite any shared assets with the same name
		for key, doc := range version.Docs {
			bookVersion.Docs[key] = doc
		}

		b.Versions[key] = *bookVersion
	}
}

type toc map[string][]TOCEntry
type versionedToc map[string]toc

func (b *Book) MakeTOC(s *Sources) error {

	entries, err := b.buildEntries(s)
	if err != nil {
		return err
	}

	//create a workdir
	b.Root.WorkDir, err = os.MkdirTemp(s.Root.SourcePath, "summary")

	//write the summary
	summary, err := os.Create(b.Root.WorkDir + "/" + SummaryFileName)
	if err != nil {
		return err
	}
	defer summary.Close()

	mg := newMarkdownGenerator()
	err = mg.GenerateSummary(entries, summary)

	return nil
}

func (b *Book) ClearWorkDir() error {
	return os.RemoveAll(b.Root.WorkDir)
}

func (b *Book) buildEntries(s *Sources) (*versionedToc, error) {
	var summary = make(versionedToc)

	shared, err := b.loadSharedEntries(s)
	if err != nil {
		return nil, err
	}

	for _, version := range s.Versions {

		configuration, err := b.buildVersionEntries(shared, version)
		if err != nil {
			return nil, err
		}
		summary[version.Version] = configuration
	}
	return &summary, nil
}

func (b *Book) loadSharedEntries(s *Sources) (toc, error) {
	file, err := os.ReadFile(s.Shared.TOC.SourcePath + "/" + s.Shared.TOC.Storage.Name())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	shared := make(toc)

	err = yaml.Unmarshal(file, &shared)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return shared, nil
}

func (b *Book) buildVersionEntries(shared toc, version Version) (toc, error) {
	//merge the shared and versioned information
	configuration := make(toc)

	b.addSharedEntries(shared, configuration)

	t, err := b.addVersionedEntries(version, configuration)
	if err != nil {
		return t, err
	}
	return configuration, nil
}

func (b *Book) addSharedEntries(shared toc, configuration toc) {
	//add the shared information for each configuration
	for key, sharedEntry := range shared {
		configuration[key] = sharedEntry
	}
}

func (b *Book) addVersionedEntries(version Version, configuration toc) (toc, error) {
	//read the versioned information
	file, err := os.ReadFile(version.TOC.SourcePath + "/" + version.TOC.Storage.Name())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	versioned := make(toc)

	err = yaml.Unmarshal(file, &versioned)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//for each toc in the sharedVersion
	for versionSection, versionTOC := range versioned {
		//iterate over our shared toc entries
		sectionExists := false
		for sharedSection, sharedTOC := range configuration {
			if versionSection == sharedSection {
				sectionExists = true
				//we already have this section and map values need to be re-assigned not changed
				var entries = make([]TOCEntry, 0)
				//add the shared entries
				for _, sharedEntry := range sharedTOC {
					entries = append(entries, sharedEntry)
				}
				//iterate over versioned entries
				for _, versionedEntry := range versionTOC {
					found := false
					for i, se := range entries {
						//if we have a matching shared entry replace it with the versioned one
						if se.Name == versionedEntry.Name {
							entries[i] = versionedEntry
							found = true
							break
						}
					}
					//it was a new entry add it
					if !found {
						entries = append(entries, versionedEntry)
					}
				}
				configuration[sharedSection] = entries
			}
		}
		//if it is not an existing section, just add it
		if !sectionExists {
			configuration[versionSection] = versionTOC
		}
	}
	return nil, nil
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
	defer func(r *os.File) {
		_ = r.Close()
	}(r)

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
