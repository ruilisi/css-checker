BIN=go
VERSION := v0.4.0
RELEASE_NOTE := "Add Path Ignores, Settings Yaml and Similarity Threshold"
.PHONY: build test-models lint revive vet

GO_PACKAGES ?= $(shell go list ./...)

git-tag:
	git tag -a $(VERSION) -m $(RELEASE_NOTE)
	git push github $(VERSION)

release: git-tag
	goreleaser release

build:
	$(BIN) build -o "dist/css-checker"

test-models:
	gotestsum --format testname --

vet:
	@echo "Running go vet..."
	${BIN} vet $(GO_PACKAGES)
	${BIN} install github.com/ruilisi/govet@v0.1.3
	${BIN} vet -vettool=$(GOPATH)/bin/govet $(GO_PACKAGES)

revive:
	GO111MODULE=on go run build/lint.go -config .revive.toml ./... || exit 1

lint: revive vet
lint-fix: fmt revive vet

coverage:
	${BIN} test -v -coverprofile cover.out .
	${BIN} tool cover -html=cover.out -o cover.html

tools:
	${BIN} install github.com/cespare/reflex@latest
	${BIN} install github.com/rakyll/gotest@latest
	${BIN} install github.com/psampaz/go-mod-outdated@latest
	${BIN} install github.com/jondot/goweight@latest
	${BIN} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	${BIN} get -t -u golang.org/x/tools/cmd/cover
	${BIN} get -t -u github.com/sonatype-nexus-community/nancy@latest
	go mod tidy

audit: tools
	${BIN} mod tidy
	${BIN} list -json -m all | nancy sleuth
