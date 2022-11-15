package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	// InfoColor ...
	InfoColor = "\033[1;34m%s\033[0m"
	// NoticeColor ...
	NoticeColor = "\033[1;36m%s\033[0m"
	// WarningColor ...
	WarningColor = "\033[1;33m%s\033[0m"
	// ErrorColor ...
	ErrorColor = "\033[1;31m%s\033[0m"
	// DebugColor ...
	DebugColor = "\033[0;36m%s\033[0m"
)

const (
	// InfoColor ...bold
	hInfoColor = "Blue"
	// NoticeColor ... bold
	hNoticeColor = "Cyan"
	// WarningColor ...bold
	hWarningColor = "Yellow"
	// ErrorColor ...bold
	hErrorColor = "Red"
	// DebugColor ...
	hDebugColor = "Cyan"
)

const (
	// Version for current version of css-checker
	Version = "0.4.1"
)

// Params setting parameters
type Params struct {
	Version             bool     `yaml:"version"`
	ColorsCheck         bool     `yaml:"colors"`
	CSS                 bool     `yaml:"css"`
	SectionsCheck       bool     `yaml:"sections"`
	SimilarityCheck     bool     `yaml:"sim"`
	SimilarityThreshold int      `yaml:"sim-threshold"`
	StyledComponents    bool     `yaml:"styled"`
	LongScriptsCheck    bool     `yaml:"long-line"`
	Path                string   `yaml:"path"`
	LongScriptLength    int      `yaml:"length-threshold"`
	Ignores             []string `yaml:"ignores"`
	Unused              bool     `yaml:"unused"`
	Unrestricted        bool     `yaml:"unrestricted"`
	ConfigPath          string   `yaml:"config"`
	ToFile				bool	 `yaml:"to-file"`
	OutputFileName      string   `yaml:"file-name"`
}

var params = Params{
	Version:          false,
	ColorsCheck:      true,
	SectionsCheck:    true,
	LongScriptsCheck: true,
	CSS:              true,
	Path:             ".",
	LongScriptLength: 20,
	Ignores:          []string{},
	ToFile:			  true,
	OutputFileName:   "css-checker.html",
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

// ScriptSummary Here, script stands for css colors and lint lines (although they are not actual scripts ^_^).
type ScriptSummary struct {
	hashValue uint64
	value     string
	scripts   []Script
	count     int
}

// SectionSummary for one section, the classes and count that occured under given paths
type SectionSummary struct {
	names []string
	value string
	count int
}

// SimilaritySummary records 2 sections' similarities and their common lines (here, duplicatedScript stands for common line)
type SimilaritySummary struct {
	sections          [2]StyleSection
	similarity        int
	duplicatedScripts []string
}

// StyleHashRecorder records sections index and original string
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

// HashOrigin hashvalue and its origin
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
	threshold := float32(params.SimilarityThreshold) / float32(100)
	for key, element := range records {
		left, right := styleList[key[0]], styleList[key[1]]
		lengthLeft, lengthRight := len(left.value), len(right.value)
		if float32(len(element)) > float32(lengthLeft)*threshold || float32(len(element)) > float32(lengthRight)*threshold {
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

// return values: (isConfigFileFound, error)
func getConf(conf *Params, path string) (bool, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return false, nil // no config file is not an error
	}

	fmt.Printf("Config YAML found, using configs in: %s\n", path)
	err = yaml.Unmarshal(buf, conf)
	if err != nil {
		fmt.Printf(ErrorColor, fmt.Sprintf("Config Error: in file %q: %v\n", path, err)) // config file in wrong format is an error
		return true, err
	}
	return true, err
}

// ParamsParse parse the given config from command line and .yaml file
func ParamsParse() {
	ignorePathsString := ""
	flag.BoolVar(&params.ColorsCheck, "colors", true, "whether to check colors")
	flag.BoolVar(&params.CSS, "css", true, "whether to check css files")
	flag.StringVar(&ignorePathsString, "ignores", "", "paths and files to be ignored (e.g. node_modules,*.example.css)")
	flag.IntVar(&params.LongScriptLength, "length-threshold", 20, "Min length of a single style value (no including the key) that to be considered as long script line")
	flag.BoolVar(&params.LongScriptsCheck, "long-line", true, "whether to check duplicated long script lines")
	flag.StringVar(&params.Path, "path", ".", "set path to files, default to be current folder")
	flag.BoolVar(&params.SectionsCheck, "sections", true, "whether to check css class duplications")
	flag.BoolVar(&params.SimilarityCheck, "sim", true, "whether to check similar css classes")
	flag.IntVar(&params.SimilarityThreshold, "sim-threshold", 80, "Threshold for Similarity Check (int only, >=20 && < 100, e.g. 80 for 80%)")
	flag.BoolVar(&params.StyledComponents, "styled", false, "checks for styled components")
	flag.BoolVar(&params.Unrestricted, "unrestricted", false, "search all files (gitignore)")
	flag.BoolVar(&params.Unused, "unused", false, "whether to check unused classes (Beta)")
	flag.BoolVar(&params.Version, "version", false, "prints current version and exits")
	flag.StringVar(&params.ConfigPath, "config", "", "set configuration file, check github.com/ruilisi/css-checker for details")
	flag.BoolVar(&params.ToFile, "to-file", true, "output result to a html file. default value is true")
	flag.StringVar(&params.OutputFileName, "file-name", "css-checker.html", "set output file name. default is css-checker.html")

	flag.Parse()
	if len(ignorePathsString) > 0 {
		params.Ignores = strings.Split(ignorePathsString, ",")
	}
	dirname, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf(ErrorColor, "Home path not found")
	}
	if strings.Contains(params.Path, "~") {
		params.Path = strings.Replace(params.Path, "~", dirname, 1) // 通过flags拿到的路径中~并不会被转译为$HOME导致读取文件错误
	}
	if strings.Contains(params.ConfigPath, "~") {
		params.ConfigPath = strings.Replace(params.ConfigPath, "~", dirname, 1)
	}
	if params.SimilarityThreshold < 20 {
		params.SimilarityThreshold = 20
	} else if params.SimilarityThreshold >= 100 {
		params.SimilarityThreshold = 99
	}
}

func main() {
	t1 := time.Now()
	ParamsParse()

	// 创建输出文件
	createOutuputFile(params.OutputFileName)
	// 是否将结果写到文件
	wtf := params.ToFile
	// 获取输出的html文件
	hf := getHtmlFile(params.OutputFileName)

	if params.Version {
		vMsg := fmt.Sprintf("<p>Version: v%s</p>\n", Version)
		// fmt.Printf("Version: v%s\n", Version)
		fmt.Printf(vMsg)
		writeToFile(hf, vMsg, wtf)
		return
	}

	// Read Config File
	configPath := params.ConfigPath
	if len(params.ConfigPath) == 0 {
		configPath = fmt.Sprintf("css-checker.yaml")
	}
	found, err := getConf(&params, configPath)
	if err != nil {
		return // config file in wrong format
	}

	// File Walk Starts
	patternsToCheck := []string{""}
	if params.StyledComponents {
		patternsToCheck = []string{"*.js", "*.jsx", "*.ts", "*.tsx"}
	}
	if params.CSS {
		patternsToCheck = append(patternsToCheck, "*.css")
	}
	files, err := WalkMatch(params.Path, WalkMatchOptions{patterns: patternsToCheck, ignores: params.Ignores, unrestricted: params.Unrestricted})
	if err != nil {
		eMsg := fmt.Sprintf("<p>No css files found at given path: %s</p>", params.Path)
		fmt.Printf(ErrorColor, fmt.Sprintf("No css files found at given path: %s", params.Path))
		writeToFile(hf, eMsg, wtf)
		return
	}
	fmt.Println("\nChecking starts. this may take seconds.")

	csf := fmt.Sprintf("<p style='color: %s'>Found %d css files. Begin to scan.</p>", hNoticeColor , len(files))
	fmt.Printf(NoticeColor, fmt.Sprintf("Found %d css files. Begin to scan.\n", len(files)))
	writeToFile(hf, csf, wtf)



	// CSS Parsing
	for _, path := range files {
		longScripts, colorScripts := SectionsParse(path, params.SimilarityThreshold)
		longScriptList = append(longScriptList, longScripts...)
		colorScriptList = append(colorScriptList, colorScripts...)
	}
	fcs := fmt.Sprintf("<p style='color: %s'>Found %d css sections. Begin to compare.</p><br/><br/>", hDebugColor, len(styleList))
	fmt.Printf(DebugColor, fmt.Sprintf("Found %d css sections. Begin to compare.\n", len(styleList)))
	writeToFile(hf, fcs, wtf)


	// Begin Checking
	dupScripts, dupColors, dupSections := []ScriptSummary{}, []ScriptSummary{}, []SectionSummary{}
	similaritySummarys := []SimilaritySummary{}
	notFoundSections := []StyleSection{}
	if params.LongScriptsCheck {
		dupScripts = DupScriptsChecker(longScriptList)
		LongScriptsWarning(dupScripts, hf, wtf)
	}
	if params.ColorsCheck {
		dupColors = DupScriptsChecker(colorScriptList)
		ColorScriptsWarning(dupColors, hf, wtf)
	}
	if params.SectionsCheck {
		dupSections = DupStyleSectionsChecker(styleList)
		StyleSectionsWarning(dupSections, hf, wtf)
	}
	if params.SimilarityCheck {
		similaritySummarys = getSimilarSections()
		SimilarSectionsWarning(similaritySummarys, params.SimilarityThreshold, hf, wtf)
	}

	if params.Unused {
		notFoundSections = UnusedClassesChecker()
		UnusedScriptsWarning(notFoundSections, hf, wtf)
	}

	t2 := time.Now()

	// Results ...
	fmt.Printf(DebugColor, fmt.Sprintf("\nCss Scan Completed.\n"))
	csc := fmt.Sprintf("<p style='color: %s'> Css Scan Completed. </p>", hDebugColor)
	writeToFile(hf, csc, wtf)
	if params.LongScriptsCheck && len(dupScripts) > 0 {
		fdl := fmt.Sprintf("<p style='color: %s'>Found %s duplicated long script values</p>", hWarningColor,  fmt.Sprintf("<span style='color: %s'>%d</span>", hErrorColor, len(dupScripts)))
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated long script values\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupScripts)))))
		writeToFile(hf, fdl, wtf)

	}
	if params.ColorsCheck && len(dupColors) > 0 {
		fdc := fmt.Sprintf("<p style='color: %s'>Found %s duplicated colors</p>", hWarningColor,fmt.Sprintf("<span style='color: %s'>%d</span>", hErrorColor, len(dupColors)))
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated colors\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupColors)))))
		writeToFile(hf, fdc, wtf)
	}
	if params.SectionsCheck && len(dupSections) > 0 {
		fdcc := fmt.Sprintf("<p style='color: %s'>Found %s duplicated css classes</p>", hWarningColor, fmt.Sprintf("<span style='color: %s'>%d</span>", hErrorColor, len(dupSections)))
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s duplicated css classes\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(dupSections)))))
		writeToFile(hf, fdcc, wtf)
	}
	if params.SimilarityCheck && len(similaritySummarys) > 0 {
		fsc := fmt.Sprintf("<p style='color: %s'>Found %s similar css classes (%d%% <= sim < 100%%)</p>", hWarningColor, fmt.Sprintf("<span style='color: %s'>%d</span>", hErrorColor, len(similaritySummarys)), params.SimilarityThreshold)
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s similar css classes (%d%% <= sim < 100%%)\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(similaritySummarys))), params.SimilarityThreshold))
		writeToFile(hf, fsc, wtf)
	}
	if params.Unused && len(notFoundSections) > 0 {
		nrfc := fmt.Sprintf("<p style='color: %s'>Found %s css classes not referred in your js/jsx/ts/tsx/htm/html code</p>", hWarningColor,  fmt.Sprintf("<span style='color: %s'>%d</span>",hErrorColor, len(notFoundSections)))
		fmt.Printf(WarningColor, fmt.Sprintf("Found %s css classes not referred in your js/jsx/ts/tsx/htm/html code\n", fmt.Sprintf(ErrorColor, fmt.Sprintf("%d", len(notFoundSections)))))
		writeToFile(hf, nrfc, wtf)
	}

	diff := t2.Sub(t1)
	if !found {
		nfy := "<p>Checking completed, you can also create a css-checker.yaml file to customize your scan. </p>"
		fmt.Println("Checking completed, you can also create a css-checker.yaml file to customize your scan.")
		writeToFile(hf, nfy, wtf)
	}
	ct := "<p>Time consumed (not including printing process): " +  shortDur(diff) + "</p>"
	fmt.Println("Time consumed (not including printing process): ", diff)
	writeToFile(hf, ct, wtf)

	writeToFile(hf, "</body></html>", wtf)
	
}


func shortDur(d time.Duration) string {
    s := d.String()
    if strings.HasSuffix(s, "m0s") {
        s = s[:len(s)-2]
    }
    if strings.HasSuffix(s, "h0m") {
        s = s[:len(s)-2]
    }
    return s
}