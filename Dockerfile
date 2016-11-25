FROM golang:1.7-alpine

RUN set -ex \
	&& apk add --no-cache \
		git \
		make


COPY ./ /go/src/github.com/cloudposse/github-authorized-keys

WORKDIR /go/src/github.com/cloudposse/github-authorized-keys

RUN make setup && go-wrapper install

ENV GITHUB_API_TOKEN
ENV GITHUB_ORGANIZATION
ENV GITHUB_TEAM
ENV GITHUB_TEAM_ID

ENV SYNC_USERS_GID
ENV SYNC_USERS_GROUPS
ENV SYNC_USERS_SHELL /bin/bash

ENTRYPOINT ["github-authorized-keys"]