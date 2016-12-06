FROM golang:1.7-alpine

COPY ./ /go/src/github.com/cloudposse/github-authorized-keys

WORKDIR /go/src/github.com/cloudposse/github-authorized-keys

ARG RUN_TESTS=0
ARG TEST_GITHUB_API_TOKEN=
ARG TEST_GITHUB_ORGANIZATION=
ARG TEST_GITHUB_TEAM=
ARG TEST_GITHUB_TEAM_ID=
ARG TEST_GITHUB_USER=
ARG TEST_ETCD_ENDPOINTS=
ARG TEST_ETCD_TTL=

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		git \
		make \
		curl \
		&& curl https://glide.sh/get | sh \
		&& make deps \
		&& if [ $RUN_TESTS -eq 1 ]; then make test ; fi \
		&& go-wrapper install \
		&& rm -rf  /go/src \
		&& apk del .build-deps

WORKDIR $GOPATH

ENV GITHUB_API_TOKEN=
ENV GITHUB_ORGANIZATION=
ENV GITHUB_TEAM=
ENV GITHUB_TEAM_ID=

ENV ETCD_ENDPOINTS=
ENV ETCD_TTL=

ENV SYNC_USERS_GID=
ENV SYNC_USERS_GROUPS=
ENV SYNC_USERS_SHELL=/bin/bash

ENTRYPOINT ["github-authorized-keys"]