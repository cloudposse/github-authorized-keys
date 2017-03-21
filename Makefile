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

SHELL = /bin/bash
export BUILD_HARNESS_PATH ?= $(shell until [ -d "build-harness" ] || [ "`pwd`" == '/' ]; do cd ..; done; pwd)/build-harness
-include $(BUILD_HARNESS_PATH)/Makefile

APP:=github-authorized-keys
COPYRIGHT_SOFTWARE:=Github Authorized Keys
COPYRIGHT_SOFTWARE_DESCRIPTION:=Use GitHub teams to manage system user accounts and authorized_keys

.PHONY : init
## Init build-harness
init:
	@curl --retry 5 --retry-delay 1 https://raw.githubusercontent.com/cloudposse/build-harness/master/bin/install.sh | bash

.PHONY : clean
## Clean build-harness
clean:
	@rm -rf $(BUILD_HARNESS_PATH)

## Build release for all supported platforms
build-all:
	gox -osarch="$(BUILD_ARCH)" -output "${RELEASE_DIR}/${APP}_{{.OS}}_{{.Arch}}"
