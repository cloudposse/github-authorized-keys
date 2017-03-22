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
