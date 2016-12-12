GO	:= $(shell which go)
GLIDE	:= $(shell which glide)
APP	?= github-authorized-keys
INSTALL_DIR ?= /usr/local/sbin

.PHONY: build
## Build binary
build: $(GO)
	$(GO) build -o $(APP)

.PHONY: test
## Run tests
test: $(GO)
  ## Find tests in all dirs expect vendor dir and dirs starts from . ##
  ## Convert local dir path to package absolute name ##
  ## Run tests ##
	find !  -path "\./\.*" ! -path "\./vendor*" -type d | \
	sed -e "s|.|github.com/cloudposse/github-authorized-keys/|" | \
	xargs -n 1 $(GO) test -v

.PHONY: deps
## Install dependencies
deps: $(GLIDE)
	$(GLIDE) update

## Clean compiled binary
clean:
	rm -f $(APP)
	$(GLIDE) remove all

## Install cli
install: $(APP)
	cp $(APP) $(INSTALL_DIR)
	chmod 555 $(INSTALL_DIR)/$(APP)

.PHONY: lint
## Lint code
lint: $(GO)
	find . ! -path "*/vendor/*" ! -path "*/.glide/*" -type f -name '*.go' | xargs -n 1 golint
	find . ! -path "*/vendor/*" ! -path "*/.glide/*" -type f -name '*.go' | xargs -n 1 $(GO) vet -v

.PHONY: fmt
fmt: $(GO)
	find . ! -path "*/vendor/*" ! -path "*/.glide/*" -type f -name '*.go' | xargs -n 1 gofmt -w -l -s

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
