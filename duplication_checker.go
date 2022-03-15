package main

import (
	"fmt"
	"strings"

	linq "github.com/ahmetb/go-linq/v3"
)

// DupStyleSectionsChecker group and count duplicated css classes
func DupStyleSectionsChecker(styleList []StyleSection) []SectionSummary {
	groups := []SectionSummary{}
	linq.From(styleList).GroupBy(func(script interface{}) interface{} {
		return script.(StyleSection).valueHash // hash value as key
	}, func(script interface{}) interface{} {
		return script
	}).Where(func(group interface{}) bool {
		return len(group.(linq.Group).Group) > 1
	}).OrderByDescending( // sort groups by its counts
		func(group interface{}) interface{} {
			return len(group.(linq.Group).Group)
		}).SelectT( // get structs out of groups
		func(group linq.Group) interface{} {
			names := []string{}
			for _, styleSection := range group.Group {
				names = append(names, fmt.Sprintf("%s << %s", styleSection.(StyleSection).name, styleSection.(StyleSection).filePath))
			}
			return SectionSummary{
				names: names,
				value: strings.Join(group.Group[0].(StyleSection).value, "\n"),
				count: len(names),
			}
		}).ToSlice(&groups)
	return groups
}

// DupScriptsChecker group and count duplicated css colors or long css lines.
func DupScriptsChecker(longScriptList []Script) []ScriptSummary {
	groups := []ScriptSummary{}
	linq.From(longScriptList).GroupBy(func(script interface{}) interface{} {
		return script.(Script).hashValue // script hashed value as key
	}, func(script interface{}) interface{} {
		return script
	}).Where(func(group interface{}) bool {
		return len(group.(linq.Group).Group) > 1
	}).OrderByDescending( // sort groups by its length
		func(group interface{}) interface{} {
			return len(group.(linq.Group).Group)
		}).SelectT( // get structs out of groups
		func(group linq.Group) interface{} {
			scripts := []Script{}
			for _, group := range group.Group {
				scripts = append(scripts, group.(Script))
			}
			return ScriptSummary{
				scripts:   scripts,
				hashValue: group.Key.(uint64),
				value:     scripts[0].value,
				count:     len(scripts),
			}
		}).ToSlice(&groups)
	return groups
}
