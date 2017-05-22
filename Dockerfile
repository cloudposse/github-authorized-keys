FROM alpine:3.5

WORKDIR /

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

RUN apk --update --no-cache add libc6-compat ca-certificates && \
    ln -s /lib /lib64

COPY ./release/github-authorized-keys_linux_amd64 /usr/bin/github-authorized-keys

ENTRYPOINT ["github-authorized-keys"]
