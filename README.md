# github-authorized-keys
Allow to provide ssh access to servers based on github teams

## Build docker image
```
docker build ./ -t github-authorized-keys
```

## Authorize command

```
docker run \                                                                      [system]
-e GITHUB_API_TOKEN={token} \
-e GITHUB_ORGANIZATION={organization} \
-e GITHUB_TEAM={team} \
github-authorized-keys authorize {user}
```

Work only for versions openssh >= 6.9



## Create users

```
docker run \                                                                      [system]
-e GITHUB_API_TOKEN={token} \
-e GITHUB_ORGANIZATION={organization} \
-e GITHUB_TEAM={team} \
-e SYNC_USERS_GID={gid OR empty} \
-e SYNC_USERS_GROUPS={comma separated groups OR empty} \
-e SYNC_USERS_SHELL={user shell} \
-v /etc:/etc \
-v /home:/home \
github-authorized-keys sync_users
```

You have to share ``/etc`` because ``adduser`` command backup ``/etc/passwd`` to  ``/etc/passwd-`` with system call
fire EXDEV error if backup are on different layers.
https://docs.docker.com/engine/userguide/storagedriver/aufs-driver/


# Demo

We use Vagrant to demonstrate how this tool works.

## Deps

Install

**Virtual box** (tested on version 4.3.26) https://www.virtualbox.org/wiki/Downloads

**Vagrant** (tested on version 1.8.4) https://www.vagrantup.com/downloads.html

**vagrant-docker-compose** plugin  with command

``vagrant plugin install vagrant-docker-compose``

## Run

Run vagrant with command


```
GITHUB_API_TOKEN={api token} \
GITHUB_ORGANIZATION={organization name} \
GITHUB_TEAM={team name} \
vagrant up
```


## Test

Login into vagrant box with command

``ssh {user}@192.168.33.10``

## Logs

You can check what is going with ssh inside vagrant box

``sudo tail -f /var/log/auth.log``