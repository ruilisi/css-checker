VERSION := v0.4.0
RELEASE_NOTE := "Add Path Ignores, Settings Yaml and Similarity Threshold"
.PHONY: build test-models

git-tag:
	git tag -a $(VERSION) -m $(RELEASE_NOTE)
	git push github $(VERSION)

release: git-tag
	goreleaser release

build:
	go build -o "css-checker"

test-models:
	gotestsum --format testname --
