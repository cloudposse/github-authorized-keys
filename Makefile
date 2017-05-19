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

#include $(shell curl -so .build-harness "https://raw.githubusercontent.com/cloudposse/build-harness/master/templates/Makefile.build-harness"; echo .build-harness)
include $(shell curl -so .build-harness "https://raw.githubusercontent.com/cloudposse/build-harness/custom-dockerfile/templates/Makefile.build-harness"; echo .build-harness)

deps:
	$(SELF) go:deps go:deps-dev go:deps-build

build:
	$(SELF) go:build
