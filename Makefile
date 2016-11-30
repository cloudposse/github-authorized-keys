GO	:= $(shell which go)
BIN	= github-authorized-keys

.PHONY: build
build: $(GO)
	$(GO) build -o $(BIN)

.PHONY: test
test: $(GO)
	$(GO) test github.com/cloudposse/github-authorized-keys/cmd


.PHONY: deps
deps: $(GO)
	$(GO) get -d -v "github.com/google/go-github/github"
	$(GO) get -d -v "golang.org/x/oauth2"
	$(GO) get -d -v "github.com/spf13/cobra/cobra"

clean:
	rm -f $(BIN)

install: $(BIN)
	cp $(BIN) /usr/local/sBIN/
	chmod 555 /usr/local/sBIN/$(BIN)


#- development targets

.PHONY: run
run: build
	./$(BIN) --config ./config.json

## Lint code
.PHONY: lint
lint: $(GO)
	golint cmd/*
	$(GO) vet -v cmd/*

## Install development dependencies
.PHONY: deps-dev
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
