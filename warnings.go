package main

import (
	"fmt"
	"strings"
	"os"
	"github.com/mazznoer/csscolorparser"
)

// SimilarSectionsWarning prints warnings for similar sections
func SimilarSectionsWarning(similaritySummarys []SimilaritySummary, sim int, f *os.File, wtf bool) {
	if len(similaritySummarys) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\n%d similar classes found as follow (%d%% <= sim < 100%%)\n.\n", len(similaritySummarys), sim))
		scf := fmt.Sprintf("<p class='t' style='color: %s'>%d similar classes found as follow (%d%% <= sim < 100%%)\n.</p>", hWarningColor, len(similaritySummarys), sim)
		writeToFile(f, scf, wtf)
		for index, summary := range similaritySummarys {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(ErrorColor, fmt.Sprintf("Sections share %d per cent similarity:\n", summary.similarity))
			ssp := fmt.Sprintf("<p style='color: %s'> <span style='color: %s'>(%d)</span> <span> Sections share %d per cent similarity: </span> </p>", hErrorColor, hWarningColor,index, summary.similarity)
			writeToFile(f, ssp, wtf)
			for _, section := range summary.sections {
				fmt.Printf(WarningColor, fmt.Sprintf("%s << %s\n", section.name, section.filePath))
				cs:= fmt.Sprintf("<p style='color: %s'> %s << %s </p>", hWarningColor, section.name, section.filePath)
				fmt.Printf(DebugColor, "\n{\n")
				cs += fmt.Sprintf("<p style='color: %s'> { </p>", hDebugColor)
				writeToFile(f, cs, wtf)
				for _, line := range section.value {
					if strings.Contains(strings.Join(summary.duplicatedScripts, "\n"), line) {
						fmt.Printf(DebugColor, fmt.Sprintln(line))
						// tc := fm.Sprintf("") 
						writeToFile(f, fmt.Sprintf("<span style='color: %s'> &nbsp;&nbsp; %s </span><br/>", hDebugColor, fmt.Sprintln(line)), wtf)
					} else {
						fmt.Println(line)
						writeToFile(f, "<p>  &nbsp;&nbsp; " + line + "</p>", wtf)
					}
				}
				fmt.Printf(DebugColor, "}\n\n")
				writeToFile(f, fmt.Sprintf("<p style='color: %s'> } </p>", hDebugColor), wtf)	
			}
		}
		fmt.Printf(WarningColor, fmt.Sprintf("For above classes, %s stands for duplicated lines\n\n\n", fmt.Sprintf(DebugColor, "Cyan Color")))
		fac := fmt.Sprintf("<p style='color: %s'> For above classes, <span style='color: %s'>Cyan Color </span> stands for duplicated lines  </p><br/><br/><br/>", hWarningColor, hDebugColor)
		writeToFile(f, fac, wtf)
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, fmt.Sprintf("No similar class found (%d%% <= sim < 100%%)\n", sim))
		writeToFile(f, fmt.Sprintf("<p style='color: %s'><span style='color: %s'>√ &nbsp;&nbsp;&nbsp;&nbsp;</span>No similar class found (%d%% <= sim < 100%%) </p>", hInfoColor,hDebugColor, sim), wtf)
	}
}

// StyleSectionsWarning prints warnings for style sections
func StyleSectionsWarning(dupStyleSections []SectionSummary,  f *os.File, wtf bool) {
	if len(dupStyleSections) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\n%d duplicated classes found as follow.\n", len(dupStyleSections)))
		dcf := fmt.Sprintf("<br/><p class='t' style='color: %s'> %d duplicated classes found as follow </p>", hWarningColor, len(dupStyleSections))
		writeToFile(f, dcf, wtf)
		ff := ""
		for index, longScript := range dupStyleSections {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(ErrorColor, fmt.Sprintf("Same class content found in %d places:\n", longScript.count))
			scc := fmt.Sprintf("<p style='color: %s'> <span style='color: %s'> (%d) </span>  Same class content found in %d places: </p>", hErrorColor, hWarningColor, index, longScript.count)
			writeToFile(f, scc, wtf)
			bf := ""
			for _, name := range longScript.names {
				fmt.Printf("\t %s\n", name)
				bf += "&nbsp;&nbsp;&nbsp;&nbsp;"+name+"<br/>"
			}
			fmt.Printf(DebugColor, fmt.Sprintf("Content:\n{\n%s\n}\n\n", longScript.value))
			bf += fmt.Sprintf("<p style='color: %s'> Content: <br/> { <br/> <pre> %s </pre> <br/> } <br/> </p>" , hDebugColor, longScript.value)
			writeToFile(f, bf, wtf)
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe content of %d duplicated contents shall be reused.\n", len(dupStyleSections)))
		ff += fmt.Sprintf("<br/><p style='color: %s'>The content of %d duplicated contents shall be reused.</p><br/>", hWarningColor, len(dupStyleSections))
		writeToFile(f, ff, wtf)
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated un-variabled color script found\n")

		writeToFile(f, fmt.Sprintf("<p style='color: %s'><span style='color: %s'>√ &nbsp;&nbsp;&nbsp;&nbsp;</span> No duplicated un-variabled color script found </p>", hInfoColor,hDebugColor), wtf)

	}
}

// ColorScriptsWarning prints warnings for unvariabled colors that used more than once
func ColorScriptsWarning(dupLongScripts []ScriptSummary, f *os.File, wtf bool) {
	if len(dupLongScripts) > 0 {
		dcf:= fmt.Sprintf("<p class='t' style='color: %s'> %d duplicated color found as follow. </p>",hWarningColor, len(dupLongScripts))
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d duplicated color found as follow.\n", len(dupLongScripts)))
		dcf += "(Colors are recommanded to be stored as variables, which can be easily updated or to be used in Themes)"
		fmt.Println("(Colors are recommanded to be stored as variables, which can be easily updated or to be used in Themes)")
		writeToFile(f, dcf, wtf)
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

			dcfinp := fmt.Sprintf(
				"<p> <span style='color: %s'>(%d)</span> <span style='color: %s'> %s ( %s )  </span> <span style='color: %s'> Found in %d places:  </span> </p>", 
				hWarningColor, 
				index, 
				hDebugColor, 
				rgbString, 
				hexString,  
				hErrorColor,
				summary.count)
			writeToFile(f, dcfinp, wtf)
			sca := ""
			for _, script := range summary.scripts {
				fmt.Printf("(%s: %s) < %s In %s\n", script.key, script.value, script.sectionName, script.filePath)
				sca += fmt.Sprintf("<p>(%s: %s) < %s In %s</p>", script.key, script.value, script.sectionName, script.filePath)
			}
			writeToFile(f, sca, wtf)
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe above %d duplicated colors shall be set to variables.\n", len(dupLongScripts)))
		tad := fmt.Sprintf("<br/><p style='color: %s'> The above %d duplicated colors shall be set to variables. </p>",hWarningColor, len(dupLongScripts))
		writeToFile(f, tad, wtf)

	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated un-variabled color script found\n")

		writeToFile(f, fmt.Sprintf("<p style='color: %s'><span style='color: %s'>√ &nbsp;&nbsp;&nbsp;&nbsp;</span> No duplicated un-variabled color script found </p>", hInfoColor,hDebugColor), wtf)

	}
}

// LongScriptsWarning prints warnings for unvariabled css long lines that used more than once
func LongScriptsWarning(dupLongScripts []ScriptSummary, f *os.File, wtf bool) {
	if len(dupLongScripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d duplicated long line found as follow.\n", len(dupLongScripts)))
		fmt.Println("(Duplicated long lines are recommanded to be extracted to variables)")
		dll := fmt.Sprintf("<p class='t' style='color: %s'>%d duplicated long line found as follow. </p>", hWarningColor, len(dupLongScripts))
		dll += "(Duplicated long lines are recommanded to be extracted to variables)"
		writeToFile(f, dll , wtf)
		for index, longScript := range dupLongScripts {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(DebugColor, longScript.value)
			fmt.Printf(ErrorColor, fmt.Sprintf(" Found in %d places:\n", longScript.count))
			lli := fmt.Sprintf(
					"<p> <span style='color: %s'> (%d) </span> <span style='color: %s'> %s </span> <span style='color: %s'> Found in %d places: </span> </p>",
					hWarningColor,
					index,
					hDebugColor,
					longScript.value,
					hErrorColor,
					longScript.count,
				)
			writeToFile(f, lli, wtf)
			ls := ""
			for _, script := range longScript.scripts {
				fmt.Printf("%s In %s\n", script.sectionName, script.filePath)
				ls += fmt.Sprintf("<p> %s In %s </p>", script.sectionName, script.filePath)
			}
			writeToFile(f, ls, wtf)
		}
		fmt.Printf(WarningColor, fmt.Sprintf("\nThe above %d duplicated long lines shall be set to variables.\n", len(dupLongScripts)))
		tad := fmt.Sprintf("<br/><p style='color: %s'> The above %d duplicated long lines shall be set to variables. </p><br/>", hWarningColor, len(dupLongScripts))
		writeToFile(f, tad, wtf)
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No duplicated long script found\n")

		writeToFile(f, fmt.Sprintf("<p style='color: %s'><span style='color: %s'>√ &nbsp;&nbsp;&nbsp;&nbsp;</span>No duplicated long script found</p>", hInfoColor,hDebugColor), wtf)

	}
}

// UnusedScriptsWarning prints warnings for unused css classes
func UnusedScriptsWarning(scripts []StyleSection, f *os.File, wtf bool) {
	if len(scripts) > 0 {
		fmt.Printf(WarningColor, fmt.Sprintf("\nOps %d classes found not used.\n", len(scripts)))
		cnu := fmt.Sprintf("<p class='t' style='color: %s'>%d classes found not used. </p>", hWarningColor, len(scripts))
		writeToFile(f, cnu, wtf)
		nc := ""
		for index, script := range scripts {
			fmt.Printf(WarningColor, fmt.Sprintf("(%d) ", index))
			fmt.Printf(DebugColor, fmt.Sprintf("%s < %s\n", script.name, script.filePath))
			nc += fmt.Sprintf(
				"<p><span style='color: %s'> (%d) </span> <span style='color: %s'>%s < %s </span>  </p>",
				hWarningColor,
				index,
				hDebugColor,
				script.name,
				script.filePath,
			)
		}
	} else {
		fmt.Printf(DebugColor, "√\t")
		fmt.Printf(InfoColor, "No unused script found\n")
		writeToFile(f, fmt.Sprintf("<p style='color: %s'><span style='color: %s'>√ &nbsp;&nbsp;&nbsp;&nbsp;</span>No unused script found</p>", hInfoColor,hDebugColor), wtf)

	}
}
