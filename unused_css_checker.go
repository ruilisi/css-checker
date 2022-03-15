package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ClassNamesSplit cases for spliting classname to parts
func ClassNamesSplit(r rune) bool {
	return r == ':' || r == '.' || r == ' ' || r == '>'
}

// JSClassNamesSplit splits js class string to separated parts
func JSClassNamesSplit(r rune) bool {
	return r == '`' || r == ' ' || r == '.' || r == '='
}

// UnusedClassesChecker checks classes that never used by any js/jsx/ts/tsx/html/htm files in the project
func UnusedClassesChecker() []StyleSection {
	files, err := WalkMatch(params.Path, WalkMatchOptions{patterns: []string{"*.js", "*.jsx", "*.ts", "*.tsx", "*.html", "*.htm"}, ignores: params.Ignores})
	notFoundSections := []StyleSection{}
	if err != nil {
		return notFoundSections
	}
	referredHashes := map[uint64]bool{}

	classReg := regexp.MustCompile(`class=["'\{]{1}[^>]*["'\}]{1}|className=["'\{]{1}[^>]*["'\}]{1}`)
	for _, filePath := range files {
		dat, err := os.ReadFile(filePath)
		if err != nil {
			return notFoundSections
		}
		result := strings.Replace(strings.Replace(string(dat), "\n", "", -1), "\r", "", -1)
		matches := classReg.FindAllStringSubmatch(result, -1)
		for _, match := range matches {
			className := strings.Replace(strings.Replace(match[0], "class", "", -1), "className", "", -1)
			className = strings.Replace(className, `"`, "", -1)
			className = strings.Replace(className, `{`, "", -1)
			className = strings.Replace(className, `}`, "", -1)
			for _, name := range strings.FieldsFunc(className, JSClassNamesSplit) {
				referredHashes[hash(name)] = true
			}
		}
	}
	for _, style := range styleList {
		names := strings.FieldsFunc(style.name, ClassNamesSplit)
		found := false
		fmt.Println(style.name, "  ", len(names), names)
		for _, name := range names {
			_, has := referredHashes[hash(name)]
			if has {
				found = true
				break
			}
		}
		if !found {
			notFoundSections = append(notFoundSections, style)
		}
	}
	return notFoundSections
}
