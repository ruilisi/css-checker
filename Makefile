VERSION := v0.3.1
RELEASE_NOTE := "Add CICD and Makefile"
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
