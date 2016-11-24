FROM golang:1.7-alpine

RUN set -ex \
	&& apk add --no-cache \
		git \
		make

COPY ./ /go/src/github.com/cloudposse/github-authorized-keys
CMD ["github-authorized-keys"]