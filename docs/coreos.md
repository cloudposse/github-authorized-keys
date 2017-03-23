# CoreOS Integration

We recommend the following defaults for deployment on CoreOS, in addition to any other variables. Stick these in your `/etc/github-authorized-keys` environment file, if using the provided `github-authorized-keys.service` in the `contrib/` folder.

```
SYNC_USERS_GID=500
SYNC_USERS_GROUPS=sudo,docker
SYNC_USERS_INTERVAL=500
SYNC_USERS_ROOT=/host
LINUX_USER_ADD_TPL=/usr/sbin/useradd --password '*' --shell {shell} {username}
LINUX_USER_ADD_WITH_GID_TPL=/usr/sbin/useradd --password '*' --shell {shell} --group {group} {username}
LINUX_USER_ADD_TO_GROUP_TPL=/usr/sbin/usermod --append --groups {group} {username}
LINUX_USER_DEL_TPL=/usr/sbin/userdel {username}
SSH_RESTART_TPL=/usr/bin/systemctl restart sshd.socket
AUTHORIZED_KEYS_COMMAND_TPL=/opt/bin/authorized-keys
```

Make sure to bind-mount your host filesystem (`/`) into `/host` on the container (e.g. `--volume /:/host`)

