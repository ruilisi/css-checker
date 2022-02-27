package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func WalkMatch(root string, patterns, ignores []string) ([]string, error) {
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
		for _, pattern := range patterns {
			if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
				return err
			} else if matched {
				matches = append(matches, path)
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}
