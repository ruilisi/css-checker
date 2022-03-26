package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicatedScriptsCheck(t *testing.T) { // same for colors and long scripts
	files := []string{"tests/normal_css.css", "tests/unformatcss.css"}
	longScriptList, colorScriptList := []Script{}, []Script{}
	for _, file := range files {
		longs, colors := SectionsParse(file, 80)
		longScriptList = append(longScriptList, longs...)
		colorScriptList = append(colorScriptList, colors...)
	}
	summaryList := DupScriptsChecker(colorScriptList)
	assert.Equal(t, len(summaryList), 1)
	assert.Equal(t, len(summaryList[0].scripts), 3)

	summaryList = DupScriptsChecker(longScriptList)
	assert.Equal(t, len(summaryList), 1)
}

func TestDuplicatedStyledComponentsCheck(t *testing.T) {
	path := "tests/sample.ts"
	longs, colors := SectionsParse(path, 80)
	summaryList := DupScriptsChecker(colors)
	assert.Equal(t, len(summaryList), 0)

	summaryList = DupScriptsChecker(longs)
	assert.Equal(t, len(summaryList), 1)
}

func TestDupStyleSectionsChecker(t *testing.T) { //same for css classes

	styleSectionFirst := StyleSection{
		name:      "firstName",
		filePath:  "tests/firstLocation",
		value:     []string{"color:'red'", "size:18"},
		valueHash: 0x5777296270491287,
	}
	styleSectionSecond := StyleSection{
		name:      "secondName",
		filePath:  "tests/secondLocation",
		value:     []string{"color:'green'", "size:18"},
		valueHash: 0xe29045cf01e7f547,
	}
	styleSectionThird := StyleSection{
		name:      "thirdName",
		filePath:  "tests/thirdLocation",
		value:     []string{"color:'red'", "size:18"},
		valueHash: 0x5777296270491287,
	}
	styleList := []StyleSection{
		styleSectionFirst,
		styleSectionSecond,
		styleSectionThird}
	sectionSummary := DupStyleSectionsChecker(styleList)
	assert.Equal(t, 2, sectionSummary[0].count)
}
