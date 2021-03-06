package main

import (
	"fmt"
	"strings"

	"github.com/mazznoer/csscolorparser"
)

// SimilarSectionsWarning prints warnings for similar sections
func SimilarSectionsWarning(similaritySummarys []SimilaritySummary, sim int) {
	if len(similaritySummarys) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\n%d similar classes found as follow (%d%% <= sim < 100%%)\n.\n", len(similaritySummarys), sim))
		for index, summary := range similaritySummarys {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(ErrorColor, fmt.Sprintf("Sections share %d per cent similarity:\n", summary.similarity))
			for _, section := range summary.sections {
				fmt.Printf(WarningColor, fmt.Sprintf("%s << %s\n", section.name, section.filePath))
				fmt.Printf(DebugColor, "\n{\n")
				for _, line := range section.value {
					if strings.Contains(strings.Join(summary.duplicatedScripts, "\n"), line) {
						fmt.Printf(DebugColor, fmt.Sprintln(line))
					} else {
						fmt.Println(line)
					}
				}
				fmt.Printf(DebugColor, "}\n\n")
			}
		}
		fmt.Printf(WarningColor, fmt.Sprintf("For above classes, %s stands for duplicated lines\n\n\n", fmt.Sprintf(DebugColor, "Cyan Color")))
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, fmt.Sprintf("No similar class found (%d%% <= sim < 100%%)\n", sim))
	}
}

// StyleSectionsWarning prints warnings for style sections
func StyleSectionsWarning(dupStyleSections []SectionSummary) {
	if len(dupStyleSections) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\n%d duplicated classes found as follow.\n", len(dupStyleSections)))
		for index, longScript := range dupStyleSections {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(ErrorColor, fmt.Sprintf("Same class content found in %d places:\n", longScript.count))
			for _, name := range longScript.names {
				fmt.Printf("\t %s\n", name)
			}
			fmt.Printf(DebugColor, fmt.Sprintf("Content:\n{\n%s\n}\n\n", longScript.value))
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe content of %d duplicated contents shall be reused.\n", len(dupStyleSections)))
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated un-variabled color script found\n")
	}
}

// ColorScriptsWarning prints warnings for unvariabled colors that used more than once
func ColorScriptsWarning(dupLongScripts []ScriptSummary) {
	if len(dupLongScripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d duplicated color found as follow.\n", len(dupLongScripts)))
		fmt.Println("(Colors are recommanded to be stored as variables, which can be easily updated or to be used in Themes)")
		for index, summary := range dupLongScripts {
			color, err := csscolorparser.Parse(summary.value)
			rgbString, hexString := summary.value, summary.value
			if err == nil {
				rgbString = color.RGBString()
				hexString = color.HexString()
			}
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(DebugColor, fmt.Sprintf("%s ( %s )", rgbString, hexString))
			fmt.Printf(ErrorColor, fmt.Sprintf(" Found in %d places:\n", summary.count))
			for _, script := range summary.scripts {
				fmt.Printf("(%s: %s) < %s In %s\n", script.key, script.value, script.sectionName, script.filePath)
			}
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe above %d duplicated colors shall be set to variables.\n", len(dupLongScripts)))
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated un-variabled color script found\n")
	}
}

// LongScriptsWarning prints warnings for unvariabled css long lines that used more than once
func LongScriptsWarning(dupLongScripts []ScriptSummary) {
	if len(dupLongScripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d duplicated long line found as follow.\n", len(dupLongScripts)))
		fmt.Println("(Duplicated long lines are recommanded to be extracted to variables)")
		for index, longScript := range dupLongScripts {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(DebugColor, longScript.value)
			fmt.Printf(ErrorColor, fmt.Sprintf(" Found in %d places:\n", longScript.count))
			for _, script := range longScript.scripts {
				fmt.Printf("%s In %s\n", script.sectionName, script.filePath)
			}
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe above %d duplicated long lines shall be set to variables.\n", len(dupLongScripts)))
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated long script found\n")
	}
}

// UnusedScriptsWarning prints warnings for unused css classes
func UnusedScriptsWarning(scripts []StyleSection) {
	if len(scripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d classes found not used.\n", len(scripts)))
		for index, script := range scripts {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(DebugColor, fmt.Sprintf("%s < %s\n", script.name, script.filePath))
		}
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No unused script found\n")
	}
}
