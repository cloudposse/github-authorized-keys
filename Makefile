APP := github-authorized-keys
COPYRIGHT_SOFTWARE := Github Authorized Keys
COPYRIGHT_SOFTWARE_DESCRIPTION := Use GitHub teams to manage system user accounts and authorized_keys

include $(shell curl -so .build-harness "https://raw.githubusercontent.com/cloudposse/build-harness/master/templates/Makefile.build-harness"; echo .build-harness)
