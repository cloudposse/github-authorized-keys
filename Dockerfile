FROM golang:1.7-alpine

COPY ./ /go/src/github.com/cloudposse/github-authorized-keys

WORKDIR /go/src/github.com/cloudposse/github-authorized-keys

ARG RUN_TESTS=0
ARG TEST_GITHUB_API_TOKEN
ARG TEST_GITHUB_ORGANIZATION
ARG TEST_GITHUB_TEAM
ARG TEST_GITHUB_TEAM_ID
ARG TEST_GITHUB_USER
ARG TEST_ETCD_ENDPOINT

# We do tests on alpine so use alpine adduser flags

ENV TEST_LINUX_USER_ADD_TPL            "adduser -D -s {shell} {username}"
ENV TEST_LINUX_USER_ADD_WITH_GID_TPL   "adduser -D -s {shell} -G {group} {username}"
ENV TEST_LINUX_USER_ADD_TO_GROUP_TPL   "adduser {username} {group}"
ENV TEST_LINUX_USER_DEL_TPL            "deluser {username}"

ENV GIN_MODE=release

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		git \
		make \
		curl \
		&& curl https://glide.sh/get | sh \
		&& make deps \
		&& ( [[ $RUN_TESTS -eq 0 ]]  ||  make test; )  \
		&& go-wrapper install \
		&& rm -rf  /go/src \
		&& apk del .build-deps

WORKDIR $GOPATH

# For production run most common user add flags
#
# We need --force-badname because github users could contains capital letters, what is not acceptable in some distributions
# Really regexp to verify badname rely on environment var that set in profile.d so we rarely hit this errors.
#
# adduser wants user name be the head and flags the tail.
ENV LINUX_USER_ADD_TPL            "adduser {username} --disabled-password --force-badname --shell {shell}"
ENV LINUX_USER_ADD_WITH_GID_TPL   "adduser {username} --disabled-password --force-badname --shell {shell} --group {group}"
ENV LINUX_USER_ADD_TO_GROUP_TPL   "adduser {username} {group}"
ENV LINUX_USER_DEL_TPL            "deluser {username}"

ENV SSH_RESTART_TPL               "/usr/sbin/service ssh force-reload"

ENV GITHUB_API_TOKEN=
ENV GITHUB_ORGANIZATION=
ENV GITHUB_TEAM=
ENV GITHUB_TEAM_ID=

ENV ETCD_ENDPOINT=
ENV ETCD_TTL=
ENV ETCD_PREFIX=/github-authorized-keys

ENV SYNC_USERS_GID=
ENV SYNC_USERS_GROUPS=
ENV SYNC_USERS_SHELL=/bin/bash
ENV SYNC_USERS_INTERVAL=

ENV INTEGRATE_SSH=false


ENV LISTEN=":301"

# For production we run container with host network, so expose is just for testing and CI\CD
EXPOSE 301

ENTRYPOINT ["github-authorized-keys"]
