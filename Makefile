# Suported platforms for compilation
BUILD_ARCH += darwin/386
BUILD_ARCH += darwin/amd64
BUILD_ARCH += linux/386
BUILD_ARCH += linux/amd64
BUILD_ARCH += linux/arm
BUILD_ARCH += linux/arm64
BUILD_ARCH += freebsd/386
BUILD_ARCH += freebsd/amd64
BUILD_ARCH += freebsd/arm
BUILD_ARCH += netbsd/386
BUILD_ARCH += netbsd/amd64
BUILD_ARCH += netbsd/arm
BUILD_ARCH += openbsd/386
BUILD_ARCH += openbsd/amd64

APP := github-authorized-keys
COPYRIGHT_SOFTWARE := Github Authorized Keys
COPYRIGHT_SOFTWARE_DESCRIPTION := Use GitHub teams to manage system user accounts and authorized_keys

include $(shell curl -so .build-harness "https://raw.githubusercontent.com/cloudposse/build-harness/master/templates/Makefile.build-harness"; echo .build-harness)

## Build release for all supported platforms
build-all:
	gox -osarch="$(BUILD_ARCH)" -output "${RELEASE_DIR}/${APP}_{{.OS}}_{{.Arch}}"
