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

build-all:
	gox -os="linux darwin freebsd openbsd netbsd plan9" -output "${RELEASE_DIR}/${APP}_{{.OS}}_{{.Arch}}"

.PHONY : clean
## Clean build-harness
clean:
	@rm -rf $(BUILD_HARNESS_PATH)





