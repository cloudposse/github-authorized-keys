# Suported platforms for compilation
RELEASE_ARCH += darwin/386
RELEASE_ARCH += darwin/amd64
RELEASE_ARCH += linux/386
RELEASE_ARCH += linux/amd64
RELEASE_ARCH += linux/arm
RELEASE_ARCH += linux/arm64
RELEASE_ARCH += freebsd/386
RELEASE_ARCH += freebsd/amd64
RELEASE_ARCH += freebsd/arm
RELEASE_ARCH += netbsd/386
RELEASE_ARCH += netbsd/amd64
RELEASE_ARCH += netbsd/arm
RELEASE_ARCH += openbsd/386
RELEASE_ARCH += openbsd/amd64

APP := github-authorized-keys
COPYRIGHT_SOFTWARE := Github Authorized Keys
COPYRIGHT_SOFTWARE_DESCRIPTION := Use GitHub teams to manage system user accounts and authorized_keys

include $(shell curl -so .build-harness "https://raw.githubusercontent.com/cloudposse/build-harness/master/templates/Makefile.build-harness"; echo .build-harness)

## Execute local deps
deps:
	$(SELF) go:deps go:deps-dev go:deps-build

## Execute local build
build:
	$(SELF) go:build

## Execute all targets
all:
	 $(SELF) go:deps-dev
	 $(SELF) go:deps-build
	 $(SELF) go:deps 
	 $(SELF) go:lint
	 $(SELF) go:test 
	 $(SELF) go:build-all

## Bring up docker compose environment
compose-up:
	docker-compose -f docker-compose-test.yaml up -d

## Entrypoint for CI
ci:
	@docker run \
		-e GIN_MODE=release \
		-e RUN_TESTS=1 \
		-e TEST_GITHUB_API_TOKEN=$(GITHUB_API_TOKEN) \
		-e TEST_GITHUB_ORGANIZATION=$(GITHUB_ORGANIZATION) \
		-e TEST_GITHUB_TEAM=$(GITHUB_TEAM) \
		-e TEST_GITHUB_TEAM_ID=$(GITHUB_TEAM_ID) \
		-e TEST_GITHUB_USER=$(GITHUB_USER) \
		-e TEST_LINUX_USER_ADD_TPL="adduser --shell {shell} {username}" \
		-e TEST_LINUX_USER_ADD_WITH_GID_TPL="adduser --shell {shell} --group {group} {username}" \
		-e TEST_LINUX_USER_ADD_TO_GROUP_TPL="adduser {username} {group}" \
		-e TEST_LINUX_USER_DEL_TPL="deluser {username}" \
		--volume=$$(pwd):/go/src/github.com/cloudposse/github-authorized-keys \
			golang:1.8 make -C /go/src/github.com/cloudposse/github-authorized-keys all
	@ls -l release/
