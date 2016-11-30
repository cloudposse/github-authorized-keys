FROM golang:1.7-alpine

COPY ./ /go/src/github.com/cloudposse/github-authorized-keys

WORKDIR /go/src/github.com/cloudposse/github-authorized-keys

RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
		git \
		make \
		&& make deps \
		&& make test \
		&& go-wrapper install \
		&& rm -rf  /go/src \
		&& apk del .build-deps

WORKDIR $GOPATH

ENV GITHUB_API_TOKEN=
ENV GITHUB_ORGANIZATION=
ENV GITHUB_TEAM=
ENV GITHUB_TEAM_ID=

ENV SYNC_USERS_GID=
ENV SYNC_USERS_GROUPS=
ENV SYNC_USERS_SHELL=/bin/bash

ENTRYPOINT ["github-authorized-keys"]