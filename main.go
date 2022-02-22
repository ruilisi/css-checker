package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	. "github.com/ahmetb/go-linq/v3"
	"github.com/aymerick/douceur/parser"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

const (
	Version = "1.0"
)

type Params struct {
	version          bool
	colorsCheck      bool
	sectionsCheck    bool
	longScriptsCheck bool
	path             string
	longScriptLength int
	ignores          []string
}

var params = Params{
	version:          false,
	colorsCheck:      true,
	sectionsCheck:    true,
	longScriptsCheck: true,
	path:             ".",
	longScriptLength: 20,
	ignores:          []string{},
}

// StyleSection ...
type StyleSection struct {
	name      string
	filePath  string
	value     []string
	valueHash uint64
}

// Script records scripts that might be extracted as variables
type Script struct {
	filePath    string
	sectionName string
	key         string
	value       string
}

// LongScriptSummary records long scripts that used more then once
type ScriptSummary struct {
	value        string
	sectionNames []string
	count        int
}

func hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func WalkMatch(root, pattern string, ignores []string) ([]string, error) {
	var matches []string
	reg := regexp.MustCompile(strings.Join(ignores, "|"))
	fmt.Println(reg)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if len(ignores) > 0 && len(reg.FindStringSubmatch(path)) > 0 {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

// SectionsParse returns StylSections, LongScripts and ColorScripts
func SectionsParse(filePath string) ([]StyleSection, []Script, []Script) {
	dat, err := os.ReadFile(filePath)
	if err != nil {
		return []StyleSection{}, []Script{}, []Script{}
	}
	stylesheet, err := parser.Parse(string(dat))
	if err != nil {
		return []StyleSection{}, []Script{}, []Script{}
	}
	styleString := strings.Replace(stylesheet.String(), "\r", "", -1)

	styleSection := StyleSection{name: "", value: []string{}, filePath: ""}
	styleList := []StyleSection{}
	longScriptList := []Script{}
	colorScriptList := []Script{}
	for _, sub := range strings.Split(styleString, "\n") {
		if strings.HasSuffix(sub, "{") {
			styleSection.name = strings.Replace(sub, "{", "", -1)
			styleSection.filePath = filePath
		} else if strings.Contains(sub, "}") {
			if len(styleSection.value) > 1 {
				sort.Strings(styleSection.value)
			}
			styleSection.valueHash = hash(strings.Join(styleSection.value, ""))
			styleList = append(styleList, styleSection)
			styleSection = StyleSection{name: "", value: []string{}, filePath: ""}
		} else {
			partials := strings.Split(sub, ": ")
			if len(partials) == 2 {
				key := strings.TrimSpace(partials[0])
				value := strings.TrimSpace(partials[1])
				// Check is Long Script
				if len(value) > params.longScriptLength && !strings.Contains(value, "var") {
					longScriptList = append(longScriptList, Script{filePath: filePath,
						sectionName: styleSection.name,
						value:       value,
						key:         key,
					})
				}
				reg := regexp.MustCompile(`#([A-Fa-f0-9]{3,6})|(rgba|hsla|rgb|hsl)\(([^}]*)\)`)
				match := reg.FindStringSubmatch(strings.ToLower(value))
				if len(match) > 0 {
					colorScriptList = append(colorScriptList, Script{
						filePath:    filePath,
						sectionName: styleSection.name,
						value:       match[0],
						key:         key,
					})
				}
			}
			styleSection.value = append(styleSection.value, sub)
		}
	}
	return styleList, longScriptList, colorScriptList
}

func DupStyleSectionsChecker(styleList []StyleSection) []ScriptSummary {
	groups := []ScriptSummary{}
	From(styleList).GroupBy(func(script interface{}) interface{} {
		return script.(StyleSection).valueHash // hash value as key
	}, func(script interface{}) interface{} {
		return script
	}).Where(func(group interface{}) bool {
		return len(group.(Group).Group) > 1
	}).OrderByDescending( // sort groups by its counts
		func(group interface{}) interface{} {
			return len(group.(Group).Group)
		}).SelectT( // get structs out of groups
		func(group Group) interface{} {
			names := []string{}
			for _, styleSection := range group.Group {
				names = append(names, fmt.Sprintf("%s << %s", styleSection.(StyleSection).name, styleSection.(StyleSection).filePath))
			}
			return ScriptSummary{
				sectionNames: names,
				value:        strings.Join(group.Group[0].(StyleSection).value, "\n"),
				count:        len(names),
			}
		}).ToSlice(&groups)
	return groups
}

func DupScriptsChecker(longScriptList []Script) []ScriptSummary {
	groups := []ScriptSummary{}
	From(longScriptList).GroupBy(func(script interface{}) interface{} {
		return script.(Script).value // script value as key
	}, func(script interface{}) interface{} {
		return fmt.Sprintf("%s: %s << %s << %s", script.(Script).key,
			script.(Script).value,
			script.(Script).sectionName,
			script.(Script).filePath,
		) // grouped info
	}).Where(func(group interface{}) bool {
		return len(group.(Group).Group) > 1
	}).OrderByDescending( // sort groups by its length
		func(group interface{}) interface{} {
			return len(group.(Group).Group)
		}).SelectT( // get structs out of groups
		func(group Group) interface{} {
			names := []string{}
			for _, name := range group.Group {
				names = append(names, name.(string))
			}
			return ScriptSummary{
				sectionNames: names,
				value:        group.Key.(string),
				count:        len(names),
			}
		}).ToSlice(&groups)
	return groups
}

func LongScriptsWarning(dupLongScripts []ScriptSummary) {
	if len(dupLongScripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d duplicated css long scripts found as follow.\n", len(dupLongScripts)))
		fmt.Println("(Duplicated long css scripts are recommanded to be extracted to variables)\n")
		for index, longScript := range dupLongScripts {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(DebugColor, longScript.value)
			fmt.Printf(ErrorColor, fmt.Sprintf(" Found in %d places:\n", longScript.count))
			for _, name := range longScript.sectionNames {
				fmt.Printf("\t %s\n", name)
			}
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe above %d duplicated css long scripts shall be set to variables.\n", len(dupLongScripts)))
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated long script found\n")
	}
}

func ColorScriptsWarning(dupLongScripts []ScriptSummary) {
	if len(dupLongScripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d duplicated color found as follow.\n", len(dupLongScripts)))
		fmt.Println("(Colors are recommanded to be stored as variables, which can be easily updated or to be used in Themes)\n")
		for index, longScript := range dupLongScripts {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(DebugColor, fmt.Sprintf("Color Value: %s", longScript.value))
			fmt.Printf(ErrorColor, fmt.Sprintf(" Found in %d places:\n", longScript.count))
			for _, name := range longScript.sectionNames {
				fmt.Printf("\t %s\n", name)
			}
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe above %d duplicated colors shall be set to variables.\n", len(dupLongScripts)))
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated un-variabled color script found\n")
	}
}

func StyleSectionsWarning(dupStyleSections []ScriptSummary) {
	if len(dupStyleSections) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\n%d duplicated css classes found as follow.\n", len(dupStyleSections)))
		for index, longScript := range dupStyleSections {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(ErrorColor, fmt.Sprintf("Same class content found in %d places:\n", longScript.count))
			for _, name := range longScript.sectionNames {
				fmt.Printf("\t %s\n", name)
			}
			fmt.Printf(DebugColor, fmt.Sprintf("Css content:\n{\n%s\n}\n\n", longScript.value))
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe content of %d duplicated css content shall be reused.\n", len(dupStyleSections)))
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated un-variabled color script found\n")
	}
}

func ParamsParse() {
	ignorePathsString := ""
	flag.BoolVar(&params.version, "version", false, "prints current version and exits")
	flag.StringVar(&params.path, "path", ".", "set path to files, default to be current folder")
	flag.StringVar(&ignorePathsString, "ignores", "", "paths and files to be ignored (e.g. node_modules,*.example.css)")
	flag.BoolVar(&params.colorsCheck, "colors", true, "whether to check colors")
	flag.BoolVar(&params.sectionsCheck, "sections", true, "whether to check sections duplications")
	flag.BoolVar(&params.longScriptsCheck, "long-line", true, "whether to check duplicated long script lines")
	flag.IntVar(&params.longScriptLength, "length-threshold", 20, "Min length of a single style value (no including the key) that to be considered as long script line")
	flag.Parse()
	if len(ignorePathsString) > 0 {
		params.ignores = strings.Split(ignorePathsString, ",")
	}
}

func main() {
	ParamsParse()
	if params.version {
		fmt.Printf("Version: v%s\n", Version)
		return
	}
	if strings.Contains(params.path, "~") {
		dirname, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf(ErrorColor, "Home path not found")
		}
		params.path = strings.Replace(params.path, "~", dirname, 1) // 通过flags拿到的路径中~并不会被转译为$HOME导致读取文件错误
	}
	files, err := WalkMatch(params.path, "*.css", params.ignores)
	if err != nil {
		fmt.Printf(ErrorColor, fmt.Sprintf("No css files found at given path: %s", params.path))
		return
	}
	fmt.Println("\nChecking starts. this may take minutes to scan.")
	fmt.Printf(NoticeColor, fmt.Sprintf("Found %d css files. Begin to scan.\n", len(files)))

	styleList := []StyleSection{}
	longScriptList := []Script{}
	colorScriptList := []Script{}
	for _, path := range files {
		list, longScripts, colorScripts := SectionsParse(path)
		styleList = append(styleList, list...)
		longScriptList = append(longScriptList, longScripts...)
		colorScriptList = append(colorScriptList, colorScripts...)
	}
	fmt.Printf(DebugColor, fmt.Sprintf("Found %d css sections. Begin to compare.\n", len(styleList)))

	dupScripts, dupColors, dupSections := []ScriptSummary{}, []ScriptSummary{}, []ScriptSummary{}
	if params.longScriptsCheck {
		dupScripts = DupScriptsChecker(longScriptList)
		LongScriptsWarning(dupScripts)
	}
	if params.colorsCheck {
		dupColors = DupScriptsChecker(colorScriptList)
		ColorScriptsWarning(dupColors)
	}
	if params.sectionsCheck {
		dupSections = DupStyleSectionsChecker(styleList)
		StyleSectionsWarning(dupSections)
	}

	fmt.Printf(DebugColor, fmt.Sprintln("Css Scan Completed."))
	if params.longScriptsCheck && len(dupScripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated long script values\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupScripts)))))
	}
	if params.colorsCheck && len(dupColors) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated colors\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupColors)))))
	}
	if params.sectionsCheck && len(dupSections) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated css classes\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupSections)))))
	}
}
