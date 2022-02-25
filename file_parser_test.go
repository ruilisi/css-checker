package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

const multiColor = "box-shadow : 0 3px 6px -4px rgb(0 0 0 / 12%), 0 6px 16px 0 #FFF, 0 9px 28px 8px rgb(0 0 0 / 5%)"

func TestColorsProcessor(t *testing.T) {
	colorReg := regexp.MustCompile(`#[A-Fa-f0-9]{3,6}|rgba\([0-9,%/ ]*\)|hsla\([0-9,%/ ]*\)|rgb\([0-9,%/ ]*\)|hsl\([0-9,%/ ]*\)`)
	colors := ColorsProcessor(multiColor, colorReg)
	assert.Equal(t, len(colors), 3)
	assert.Equal(t, colors[0].original, "rgb(0 0 0 / 12%)")
	assert.Equal(t, colors[1].rgb, "rgb(255,255,255)")
}
func TestClassSectionsProcessor(t *testing.T) {
	files := []string{"tests/normal_css.css", "tests/unformatcss.css"}
	colorCounts := []int{3, 1}
	sectionsCounts := []int{4, 7}
	styleList = []StyleSection{}
	for index, file := range files {
		longScriptList, colorScriptList = SectionsParse(file)
		assert.Equal(t, len(longScriptList), 1)
		assert.Equal(t, longScriptList[0].key, "transition")
		assert.Equal(t, len(colorScriptList), colorCounts[index])
		assert.Equal(t, len(styleList), sectionsCounts[index])
	}
}
