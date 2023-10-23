package pages

import (
	"testing"
)

func TestToSectionSort(t *testing.T) {
	//arrange
	var book = make(VersionedToc)

	book["9"] = makeVersion9()
	book["10"] = makeVersion10()

	//act
	orderedVersionTocs := book.Sort()

	checkVersion9(t, orderedVersionTocs[0])
	checkVersion10(t, orderedVersionTocs[1])
}

func checkVersion9(t *testing.T, orderedVersionToc OrderedVersionTocs) {

	if orderedVersionToc.Version != "9" {
		t.Errorf("Expected 9 to be first, got %s", orderedVersionToc.Version)
	}

	if orderedVersionToc.Order != 9 {
		t.Errorf("Expected 9 to be first, got %d", orderedVersionToc.Order)
	}

	//assert
	if orderedVersionToc.Sections[0].Name != "SectionOne" {
		t.Errorf("Expected SectionOne to be first, got %s", orderedVersionToc.Sections[0].Name)
	} else {
		if orderedVersionToc.Sections[0].Section.Entries[0].Name != "DocumentFour" {
			t.Errorf("Expected DocumentFour to be first, got %s", orderedVersionToc.Sections[0].Section.Entries[0].Name)
		}
		if orderedVersionToc.Sections[0].Section.Entries[1].Name != "DocumentThree" {
			t.Errorf("Expected DocumentThree to be second, got %s", orderedVersionToc.Sections[0].Section.Entries[1].Name)
		}
	}

	if orderedVersionToc.Sections[1].Name != "SectionTwo" {
		t.Errorf("Expected SectionTwo to be second, got %s", orderedVersionToc.Sections[1].Name)
	} else {
		if orderedVersionToc.Sections[1].Section.Entries[0].Name != "DocumentTwo" {
			t.Errorf("Expected DocumentTwo to be first, got %s", orderedVersionToc.Sections[1].Section.Entries[0].Name)
		}
		if orderedVersionToc.Sections[1].Section.Entries[1].Name != "DocumentOne" {
			t.Errorf("Expected DocumentOne to be second, got %s", orderedVersionToc.Sections[1].Section.Entries[1].Name)
		}
	}

	if orderedVersionToc.Sections[2].Name != "SectionThree" {
		t.Errorf("Expected SectionThree to be third, got %s", orderedVersionToc.Sections[2].Name)
	} else {
		if orderedVersionToc.Sections[2].Section.Entries[0].Name != "DocumentFive" {
			t.Errorf("Expected DocumentFive to be first, got %s", orderedVersionToc.Sections[2].Section.Entries[0].Name)
		}
		if orderedVersionToc.Sections[2].Section.Entries[1].Name != "DocumentSix" {
			t.Errorf("Expected DocumentSix to be second, got %s", orderedVersionToc.Sections[2].Section.Entries[1].Name)
		}
	}
}

func checkVersion10(t *testing.T, orderedVersionToc OrderedVersionTocs) {

	if orderedVersionToc.Version != "10" {
		t.Errorf("Expected 10 to be first, got %s", orderedVersionToc.Version)
	}

	if orderedVersionToc.Order != 10 {
		t.Errorf("Expected 10 to be first, got %d", orderedVersionToc.Order)
	}

	//assert
	if orderedVersionToc.Sections[0].Name != "SectionTwo" {
		t.Errorf("Expected SectionTwo to be first, got %s", orderedVersionToc.Sections[0].Name)
	} else {
		if orderedVersionToc.Sections[0].Section.Entries[0].Name != "DocumentTwo" {
			t.Errorf("Expected DocumentTwo to be first, got %s", orderedVersionToc.Sections[0].Section.Entries[0].Name)
		}
		if orderedVersionToc.Sections[0].Section.Entries[1].Name != "DocumentOne" {
			t.Errorf("Expected DocumentOne to be second, got %s", orderedVersionToc.Sections[0].Section.Entries[1].Name)
		}
	}

	if orderedVersionToc.Sections[1].Name != "SectionThree" {
		t.Errorf("Expected SectionThree to be second, got %s", orderedVersionToc.Sections[1].Name)
	} else {
		if orderedVersionToc.Sections[1].Section.Entries[0].Name != "DocumentFive" {
			t.Errorf("Expected DocumentFive to be first, got %s", orderedVersionToc.Sections[1].Section.Entries[0].Name)
		}
		if orderedVersionToc.Sections[1].Section.Entries[1].Name != "DocumentSix" {
			t.Errorf("Expected DocumentSix to be second, got %s", orderedVersionToc.Sections[1].Section.Entries[1].Name)
		}
	}
}

func makeVersion9() (toc Toc) {

	toc = make(Toc)

	toc["SectionTwo"] = TOCSection{
		Order: 10,
		Entries: []TOCEntry{
			{
				Name:  "DocumentOne",
				File:  "DocumentOne.md",
				Order: 2,
			},
			{
				Name:  "DocumentTwo",
				File:  "DocumentTwo.md",
				Order: 1,
			},
		},
	}

	toc["SectionThree"] = TOCSection{
		Order: 15,
		Entries: []TOCEntry{
			{
				Name:  "DocumentSix",
				File:  "DocumentSix.md",
				Order: 18,
			},
			{
				Name:  "DocumentFive",
				File:  "DocumentFive.md",
				Order: 9,
			},
		},
	}

	toc["SectionOne"] = TOCSection{
		Order: 5,
		Entries: []TOCEntry{
			{
				Name:  "DocumentThree",
				File:  "DocumentThree.md",
				Order: 16,
			},
			{
				Name:  "DocumentFour",
				File:  "DocumentFour.md",
				Order: 12,
			},
		},
	}
	return toc
}

func makeVersion10() (toc Toc) {

	toc = make(Toc)

	toc["SectionTwo"] = TOCSection{
		Order: 10,
		Entries: []TOCEntry{
			{
				Name:  "DocumentOne",
				File:  "DocumentOne.md",
				Order: 2,
			},
			{
				Name:  "DocumentTwo",
				File:  "DocumentTwo.md",
				Order: 1,
			},
		},
	}

	toc["SectionThree"] = TOCSection{
		Order: 15,
		Entries: []TOCEntry{
			{
				Name:  "DocumentSix",
				File:  "DocumentSix.md",
				Order: 18,
			},
			{
				Name:  "DocumentFive",
				File:  "DocumentFive.md",
				Order: 9,
			},
		},
	}

	return toc
}
