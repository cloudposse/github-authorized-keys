FROM golang:1.7-alpine

COPY ./ /go/src/github.com/cloudposse/github-authorized-keys

WORKDIR /go/src/github.com/cloudposse/github-authorized-keys

ARG RUN_TESTS=0
ARG TEST_GITHUB_API_TOKEN=
ARG TEST_GITHUB_ORGANIZATION=
ARG TEST_GITHUB_TEAM=
ARG TEST_GITHUB_TEAM_ID=
ARG TEST_GITHUB_USER=

# We do tests on alpine so use alpine adduser flags

ENV TEST_LINUX_USER_ADD_TPL            "adduser -D -s {shell} {username}"
ENV TEST_LINUX_USER_ADD_WITH_GID_TPL   "adduser -D -s {shell} -G {group} {username}"
ENV TEST_LINUX_USER_ADD_TO_GROUP_TPL   "adduser {username} {group}"
ENV TEST_LINUX_USER_DEL_TPL            "deluser {username}"

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		git \
		make \
		&& make deps \
		&& if [ $RUN_TESTS -eq 1 ]; then make deps-dev && make test ; fi \
		&& go-wrapper install \
		&& rm -rf  /go/src \
		&& apk del .build-deps

WORKDIR $GOPATH

# For production run most common user add flags

ENV LINUX_USER_ADD_TPL            "adduser --disabled-password  --gecos '' --shell {shell} {username}"
ENV LINUX_USER_ADD_WITH_GID_TPL   "adduser --disabled-password  --gecos '' --shell {shell} --group {group} {username}"
ENV LINUX_USER_ADD_TO_GROUP_TPL   "adduser {username} {group}"
ENV LINUX_USER_DEL_TPL            "deluser {username}"

ENV GITHUB_API_TOKEN=
ENV GITHUB_ORGANIZATION=
ENV GITHUB_TEAM=
ENV GITHUB_TEAM_ID=

ENV SYNC_USERS_GID=
ENV SYNC_USERS_GROUPS=
ENV SYNC_USERS_SHELL=/bin/bash

ENTRYPOINT ["github-authorized-keys"]