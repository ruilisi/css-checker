package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicatedScriptsCheck(t *testing.T) { // same for colors and long scripts
	files := []string{"tests/normal_css.css", "tests/unformatcss.css"}
	longScriptList, colorScriptList := []Script{}, []Script{}
	for _, file := range files {
		longs, colors := SectionsParse(file)
		longScriptList = append(longScriptList, longs...)
		colorScriptList = append(colorScriptList, colors...)
	}
	fmt.Println(colorScriptList)
	summaryList := DupScriptsChecker(colorScriptList)
	assert.Equal(t, len(summaryList), 1)
	assert.Equal(t, len(summaryList[0].scripts), 3)

	summaryList = DupScriptsChecker(longScriptList)
	assert.Equal(t, len(summaryList), 1)
}
