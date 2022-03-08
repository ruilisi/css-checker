package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	gitignore "github.com/iriri/minimal/gitignore"
)

type WalkMatchOptions struct {
	ignores      []string
	patterns     []string
	unrestricted bool
}

func WalkMatch(root string, options WalkMatchOptions) ([]string, error) {
	var matches []string
	ignores := options.ignores
	from, ignoreListErr := gitignore.From(fmt.Sprintf("%s/.gitignore", root))
	if options.unrestricted {
		from, ignoreListErr = gitignore.New()
	}
	reg := regexp.MustCompile(strings.Join(ignores, "|"))
	err := from.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if len(options.ignores) > 0 && len(reg.FindStringSubmatch(path)) > 0 {
			return nil
		}
		if !options.unrestricted && ignoreListErr == nil {
			if from.Match(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		for _, pattern := range options.patterns {
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
