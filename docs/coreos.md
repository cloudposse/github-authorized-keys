# CoreOS Integration

We recommend the following defaults for deployment on CoreOS.

```
SYNC_USERS_GID=500
SYNC_USERS_GROUPS=sudo,docker
SYNC_USERS_INTERVAL=500
SYNC_USERS_ROOT=/host
LINUX_USER_ADD_TPL=useradd --password '*' --shell {shell} {username}
LINUX_USER_ADD_WITH_GID_TPL=adduser --password '*' --shell {shell} --group {group} {username}
SSH_RESTART_TPL=/usr/bin/systemctl restart sshd.socket
```

Make sure to bind-mount your host filesystem (`/`) into `/host` on the container (e.g. `--volume /:/host`)

