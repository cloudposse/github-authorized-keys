GO	:= $(shell which go)
BIN	= github-authorized-keys

.PHONY: build
## Build binary
build: $(GO)
	$(GO) build -o $(BIN)

.PHONY: test
## Run tests
test: $(GO)
	$(GO) test github.com/cloudposse/github-authorized-keys/cmd


.PHONY: deps
## Install dependencies
deps: $(GO)
	$(GO) get -d -v "github.com/google/go-github/github"
	$(GO) get -d -v "golang.org/x/oauth2"
	$(GO) get -d -v "github.com/spf13/cobra/cobra"

## Clean compiled binary
clean:
	rm -f $(BIN)

## Install cli
install: $(BIN)
	cp $(BIN) /usr/local/sbin/
	chmod 555 /usr/local/sbin/$(BIN)

.PHONY: lint
## Lint code
lint: $(GO)
	golint cmd/*
	golint *.go
	$(GO) vet -v cmd/*
	$(GO) vet -v *.go

.PHONY: deps-dev
## Install development dependencies
deps-dev: $(GO)
	$(GO) get -d -v "github.com/golang/lint"
	$(GO) install -v "github.com/golang/lint/golint"

## This help screen
help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+:/ { \
	  helpMessage = match(lastLine, /^## (.*)/); \
	  if (helpMessage) { \
	    helpCommand = $$1; \
	    helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
      gsub("\\\\", "", helpCommand); \
      gsub(":+$$", "", helpCommand); \
	    printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
	  } \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"
