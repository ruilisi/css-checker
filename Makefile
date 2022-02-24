VERSION := v0.2.0
RELEASE_NOTE := "Add CICD and Makefile"
.PHONY: build

git-tag:
	git tag -a $(VERSION) -m $(RELEASE_NOTE)
	git push github $(VERSION)

release: git-tag
	goreleaser release

build:
	go build -o "css-checker"
