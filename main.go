package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strings"
	"time"
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
	similarityCheck  bool
	longScriptsCheck bool
	path             string
	longScriptLength int
	ignores          []string
	unused           bool
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
	hashValue   uint64
}

type ScriptSummary struct {
	hashValue uint64
	value     string
	scripts   []Script
	count     int
}

type SectionSummary struct {
	names []string
	value string
	count int
}

type SimilaritySummary struct {
	sections          [2]StyleSection
	similarity        int
	duplicatedScripts []string
}

type StyleHashRecorder struct {
	sectionIndex int
	originString string
}

var styleList = []StyleSection{}
var longScriptList = []Script{}
var colorScriptList = []Script{}

var hashCounters = map[uint64][]StyleHashRecorder{} // hashValue -> section

func hash(s string) uint64 {
	h := fnv.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

type HashOrigin struct {
	hash   uint64
	origin string
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getSimilarSections() []SimilaritySummary {
	records := map[[2]int][]HashOrigin{}
	summary := []SimilaritySummary{}

	// Convert map LineHash -> Section => [SectionIndex1][SectionIndex2] <-> Duplicated Hashes [O(n)], n for identical hash, section stands for css class
	for key, element := range hashCounters {
		if len(element) < 2 {
			continue
		}
		for i := 0; i < len(element)-1; i++ {
			for j := i + 1; j < len(element); j++ {
				if element[i].sectionIndex < element[j].sectionIndex {
					if record, found := records[[2]int{element[i].sectionIndex, element[j].sectionIndex}]; found {
						records[[2]int{element[i].sectionIndex, element[j].sectionIndex}] = append(record, HashOrigin{hash: key, origin: element[i].originString})
					} else {
						records[[2]int{element[i].sectionIndex, element[j].sectionIndex}] = []HashOrigin{{hash: key, origin: element[i].originString}}
					}
				}
			}
		}
	}

	// In map: [SectionIndex1][SectionIndex2] -> Duplicated Hashes, number of the duplicated hashes stands for duplicated lines between classes.
	for key, element := range records {
		left, right := styleList[key[0]], styleList[key[1]]
		lengthLeft, lengthRight := len(left.value), len(right.value)
		if float32(len(element)) > float32(lengthLeft)*0.8 || float32(len(element)) > float32(lengthRight)*0.8 {
			if len(element) == min(lengthLeft, lengthRight) {
				continue
			}
			duplicatedStrings := []string{}
			for _, hashOrigin := range element {
				duplicatedStrings = append(duplicatedStrings, hashOrigin.origin)
			}
			summary = append(summary, SimilaritySummary{
				sections:          [2]StyleSection{left, right},
				similarity:        100 * len(element) / min(lengthLeft, lengthRight),
				duplicatedScripts: duplicatedStrings,
			})
		}
	}
	sort.SliceStable(summary, func(i, j int) bool {
		return summary[i].similarity < summary[j].similarity
	})
	return summary
}

func ParamsParse() {
	ignorePathsString := ""
	flag.BoolVar(&params.version, "version", false, "prints current version and exits")
	flag.StringVar(&params.path, "path", ".", "set path to files, default to be current folder")
	flag.StringVar(&ignorePathsString, "ignores", "", "paths and files to be ignored (e.g. node_modules,*.example.css)")
	flag.BoolVar(&params.colorsCheck, "colors", true, "whether to check colors")
	flag.BoolVar(&params.sectionsCheck, "sections", true, "whether to check css class duplications")
	flag.BoolVar(&params.similarityCheck, "sim", true, "whether to check similar css classes (>=80% && < 100%)")
	flag.BoolVar(&params.longScriptsCheck, "long-line", true, "whether to check duplicated long script lines")
	flag.BoolVar(&params.unused, "unused", false, "whether to check unused classes (Beta)")
	flag.IntVar(&params.longScriptLength, "length-threshold", 20, "Min length of a single style value (no including the key) that to be considered as long script line")
	flag.Parse()
	if len(ignorePathsString) > 0 {
		params.ignores = strings.Split(ignorePathsString, ",")
	}
}

func main() {
	t1 := time.Now()
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
	files, err := WalkMatch(params.path, []string{"*.css"}, params.ignores)
	if err != nil {
		fmt.Printf(ErrorColor, fmt.Sprintf("No css files found at given path: %s", params.path))
		return
	}
	fmt.Println("\nChecking starts. this may take minutes to scan.")
	fmt.Printf(NoticeColor, fmt.Sprintf("Found %d css files. Begin to scan.\n", len(files)))

	for _, path := range files {
		longScripts, colorScripts := SectionsParse(path)
		longScriptList = append(longScriptList, longScripts...)
		colorScriptList = append(colorScriptList, colorScripts...)
	}
	fmt.Printf(DebugColor, fmt.Sprintf("Found %d css sections. Begin to compare.\n", len(styleList)))

	dupScripts, dupColors, dupSections := []ScriptSummary{}, []ScriptSummary{}, []SectionSummary{}
	similaritySummarys := []SimilaritySummary{}
	notFoundSections := []StyleSection{}
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
	if params.similarityCheck {
		similaritySummarys = getSimilarSections()
		SimilarSectionsWarning(similaritySummarys)
	}

	if params.unused {
		notFoundSections = UnusedClassesChecker()
		UnusedScriptsWarning(notFoundSections)
	}

	t2 := time.Now()

	fmt.Printf(DebugColor, fmt.Sprintf("\nCss Scan Completed.\n"))
	if params.longScriptsCheck && len(dupScripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated long script values\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupScripts)))))
	}
	if params.colorsCheck && len(dupColors) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated colors\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupColors)))))
	}
	if params.sectionsCheck && len(dupSections) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated css classes\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupSections)))))
	}
	if params.similarityCheck && len(dupSections) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s similar css classes\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(similaritySummarys)))))
	}
	if params.unused && len(notFoundSections) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s css classes not referred in your js/jsx/ts/tsx/htm/html code\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(notFoundSections)))))
	}

	diff := t2.Sub(t1)
	fmt.Println("Time consumed (not including printing process): ", diff)
}
