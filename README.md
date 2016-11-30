# github-authorized-keys
Allow to provide ssh access to servers based on github teams

## Getting started

Tool is writen on go lang and provide command line interface.
It is possible to run this command as simple cli application or in docker container.

### Run as cli

#### Requirements

  Go lang 1.7.x

#### Install

  Compile with command
  ```
make deps
make build
make install
  ```

After installation you could find command as
```
/usr/local/sbin/github-authorized-keys
```

#### Run

##### Authorize command
You can specify params as flags

```
/usr/local/sbin/github-authorized-keys \
--token={token} \
--org={organization} \
--team={team} \
authorize {user}
```

or as environment variables


```
GITHUB_API_TOKEN={token} \
GITHUB_ORGANIZATION={organization} \
GITHUB_TEAM={team} \
/usr/local/sbin/github-authorized-keys authorize {user}
```

or you can mix that approaches

##### Create users

You can specify params as flags

```
/usr/local/sbin/github-authorized-keys \
--token={token} \
--org={organization} \
--team={team} \
--gid={user gid} \
--groups={comma separated secondary groups names} \
--shell={user shell} \
sync-users
```

or as environment variables


```
GITHUB_API_TOKEN={token} \
GITHUB_ORGANIZATION={organization} \
GITHUB_TEAM={team} \
SYNC_USERS_GID={gid OR empty} \
SYNC_USERS_GROUPS={comma separated groups OR empty} \
SYNC_USERS_SHELL={user shell} \
/usr/local/sbin/github-authorized-keys sync-users
```

or you can mix that approaches

### Run in containers

#### Requirements

  Docker

#### Build docker image
```
docker build ./ -t github-authorized-keys
```

### Run docker image

#### Authorize command

```
docker run \
-e GITHUB_API_TOKEN={token} \
-e GITHUB_ORGANIZATION={organization} \
-e GITHUB_TEAM={team} \
github-authorized-keys authorize {user}
```


#### Create users

```
docker run \
-e GITHUB_API_TOKEN={token} \
-e GITHUB_ORGANIZATION={organization} \
-e GITHUB_TEAM={team} \
-e SYNC_USERS_GID={gid OR empty} \
-e SYNC_USERS_GROUPS={comma separated groups OR empty} \
-e SYNC_USERS_SHELL={user shell} \
-v /etc:/etc \
-v /home:/home \
github-authorized-keys sync-users
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