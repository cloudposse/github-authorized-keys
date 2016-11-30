GO	:= $(shell which go)
bin	= github-authorized-keys

.PHONY: build
build: $(GO)
	$(GO) build -o $(bin)

.PHONY: test
test: $(GO)
	$(GO) test github.com/cloudposse/github-authorized-keys/cmd


.PHONY: setup
setup: $(GO)
	$(GO) get -d -v "github.com/google/go-github/github"
	$(GO) get -d -v "golang.org/x/oauth2"
	$(GO) get -d -v "github.com/spf13/cobra/cobra"

clean:
	rm -f $(bin)

install: $(bin)
	cp $(bin) /usr/local/sbin/
	chmod 555 /usr/local/sbin/$(bin)

#- development targets

.PHONY: run
run: build
	./$(bin) --config ./config.json


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