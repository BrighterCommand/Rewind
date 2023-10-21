package pages

import (
	"testing"
)

func TestToSectionSort(t *testing.T) {
	//arrange

	var toc = make(Toc)
	toc["SectionTwo"] = TOCSection{
		Order: 10,
		Entries: []TOCEntry{
			TOCEntry{
				Name:  "DocumentOne",
				File:  "DocumentOne.md",
				Order: 2,
			},
			TOCEntry{
				Name:  "DocumentTwo",
				File:  "DocumentTwo.md",
				Order: 1,
			},
		},
	}

	toc["SectionThree"] = TOCSection{
		Order: 15,
		Entries: []TOCEntry{
			TOCEntry{
				Name:  "DocumentSix",
				File:  "DocumentSix.md",
				Order: 18,
			},
			TOCEntry{
				Name:  "DocumentFive",
				File:  "DocumentFive.md",
				Order: 9,
			},
		},
	}

	toc["SectionOne"] = TOCSection{
		Order: 5,
		Entries: []TOCEntry{
			TOCEntry{
				Name:  "DocumentThree",
				File:  "DocumentThree.md",
				Order: 16,
			},
			TOCEntry{
				Name:  "DocumentFour",
				File:  "DocumentFour.md",
				Order: 12,
			},
		},
	}

	//act
	sortedTOC := toc.Sort()

	//assert
	if sortedTOC[0].Name != "SectionOne" {
		t.Errorf("Expected SectionOne to be first, got %s", sortedTOC[0].Name)
	}

	if sortedTOC[1].Name != "SectionTwo" {
		t.Errorf("Expected SectionTwo to be second, got %s", sortedTOC[1].Name)
	}

	if sortedTOC[2].Name != "SectionThree" {
		t.Errorf("Expected SectionThree to be third, got %s", sortedTOC[2].Name)
	}

	//ToDo: Check the entries within the section are sorted too

}
