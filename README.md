# Github Authorized Keys
Allow to provide ssh access to servers based on github teams

----

## Demo

We use Vagrant to demonstrate how this tool works.

### Deps

Install

**[Virtual box](https://www.virtualbox.org/wiki/Downloads)** (tested on version 4.3.26)

**[Vagrant](https://www.vagrantup.com/downloads.html)** (tested on version 1.8.4)

**[vagrant-docker-compose](https://github.com/leighmcculloch/vagrant-docker-compose)** plugin
  with command

``vagrant plugin install vagrant-docker-compose``

### Run

Vagrant up required some configs.
There are 2 cases to do this.

#### With config file

Copy .github-authorized-keys-demo.default.yml to .github-authorized-keys-demo.yml

```
cp .github-authorized-keys-demo.default.yml .github-authorized-keys-demo.yml
```

and set required values in that file.

Then you can simple run

```
vagrant up
```


#### With environment variables

Run vagrant with command


```
GITHUB_API_TOKEN={api token} \
GITHUB_ORGANIZATION={organization name} \
GITHUB_TEAM={team name} \
vagrant up
```


### Test

Login into vagrant box with command

``ssh -o "UserKnownHostsFile /dev/null" {github username}@192.168.33.10``

### Logs

You can check what is going with ssh inside vagrant box

``sudo tail -f /var/log/auth.log``

----------------

## Getting started

Use GitHub teams to manage system user accounts and authorized_keys.
This tool consists of two parts
1. REST API that provide public keys for users (required for sshd_config AuthorizedKeysCommand)
2. Internal cron job schedule task to create users on linux machine

### Requirements

#### Use as deamon

  * [Go lang 1.7.x](https://golang.org/)
  * [glide](https://github.com/Masterminds/glide)

#### Use in containers

  * [Docker](https://docs.docker.com/engine/installation)

### Install

#### Use as daemon


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

#### Use in containers

  Build docker image

 ```
 docker build ./ -t github-authorized-keys
 ```

### Start

##### Use as daemon

To start daemon run cli command with configuration params

You can specify params as flags

```
/usr/local/sbin/github-authorized-keys \
--github-api-token={token} \
--github-organization={organization} \
--github-team={team} \
--sync-users-gid={user gid} \
--sync-users-groups={comma separated secondary groups names} \
--sync-users-shell={user shell} \
--sync-users-root={root directory - default "/"} \
--sync-users-interval={seconds - default 300} \
--etcdctl-endpoint={etcd endpoints comma separeted - optional} \
--etcdctl-ttl={etcd ttl - default 1 day} \
--etcdctl-prefix={prefix or path to store data - default /github-authorized-keys}
```

or as environment variables


```
GITHUB_API_TOKEN={token} \
GITHUB_ORGANIZATION={organization} \
GITHUB_TEAM={team} \
SYNC_USERS_GID={gid OR empty} \
SYNC_USERS_GROUPS={comma separated groups OR empty} \
SYNC_USERS_SHELL={user shell} \
SYNC_USERS_ROOT={root directory - default "/"} \
SYNC_USERS_INTERVAL={seconds - default 300} \
ETCDCTL_ENDPOINT={etcd endpoints comma separeted - optional} \
ETCDCTL_TTL={etcd ttl - default 1 day} \
ETCDCTL_PREFIX={prefix or path to store data - default /github-authorized-keys} \
/usr/local/sbin/github-authorized-keys authorize {user}
```

or you can mix that approaches

##### Use in containers

You can specify params  as environment variables

```
docker run \
-v /:/{root directory} \
--expose "301:301"
-e GITHUB_API_TOKEN={token} \
-e GITHUB_ORGANIZATION={organization} \
-e GITHUB_TEAM={team} \
-e SYNC_USERS_GID={gid OR empty} \
-e SYNC_USERS_GROUPS={comma separated groups OR empty} \
-e SYNC_USERS_SHELL={user shell} \
-e SYNC_USERS_ROOT={root directory} \
-e SYNC_USERS_INTERVAL={seconds - default 300} \
-e ETCDCTL_ENDPOINT={etcd endpoints comma separeted - optional} \
-e ETCDCTL_TTL={etcd ttl - default 1 day} \
-e ETCDCTL_PREFIX={prefix or path to store data - default /github-authorized-keys} \
github-authorized-keys
```

or as flags

```
docker run \
-v /:/{root directory} \
--expose "301:301"
github-authorized-keys
--github-api-token={token} \
--github-organization={organization} \
--github-team={team} \
--sync-users-gid={user gid} \
--sync-users-groups={comma separated secondary groups names} \
--sync-users-shell={user shell} \
--sync-users-root={root directory - default "/"} \
--sync-users-interval={seconds - default 300} \
--etcdctl-endpoint={etcd endpoints comma separeted - optional} \
--etcdctl-ttl={etcd ttl - default 1 day} \
--etcdctl-prefix={prefix or path to store data - default /github-authorized-keys}
```
or you can mix that approaches

### Usage

#### Authorize

To integrate system with ssh you need to config AuthorizedKeysCommand in sshd_config to call authorize command (using shell wrapper).

In case of running in container you have to share host ``/`` into ``/{root directory}`` because  ``adduser`` command could differs on different Linux distribs and we need to use host one.
Also that means you need to specify  sync-users-root param to point to that directory.


Authorization REST API use ETCD to temporary cache user's public keys.
If github.com is not available command fallback to ETCD storage.

ETCD endpoints param is optional, if not specify caching and fallback disabled.


#### Create users

Linux users will be synchronized according to team members every 1 second

In case of running in container you have to share host ``/`` into ``/{root directory}`` because  ``adduser`` command could differs on different Linux distribs and we need to use host one.
Also that means you need to specify  sync-users-root param to point to that directory.

##### Templating commands

 Command sync-users rely on OS commands to mange users. We use templates for commands.
 Templates could be overridden with environment variables.

 Following templates are used:

###### Add user

**Default template:**

  ```
adduser {username} --disabled-password --force-badname --shell {shell}
  ```

**Valid placeholders:**

1. **_{username}_** - User login name
2. **_{shell}_**    - User shell

**Environment variable:**

LINUX_USER_ADD_TPL


###### Add user with primary group

**Default template:**

  ```
adduser {username} --disabled-password --force-badname --shell {shell} --group {group}
  ```

**Valid placeholders:**

1. **_{username}_** - User login name
2. **_{shell}_**    - User shell
3. **_{group}_**    - User primary group name
4. **_{gid}_**      - User primary group id

**Environment variable:**

LINUX_USER_ADD_WITH_GID_TPL

###### Add user to secondary group

**Default template:**

  ```
adduser {username} {group}
  ```

**Valid placeholders:**

1. **_{username}_** - User login name
2. **_{group}_**    - User primary group name

**Environment variable:**

LINUX_USER_ADD_TO_GROUP_TPL

###### Delete user

**Template:**

  ```
deluser {username}
  ```

**Valid placeholders:**

1. **_{username}_** - User login name

**Environment variable:**

LINUX_USER_DEL_TPL

------------

## Development

### Requirements

  * [Go lang 1.7.x](https://golang.org/)
  * [glide](https://github.com/Masterminds/glide)
  * [Make](https://en.wikipedia.org/wiki/Make_(software))
  * [Docker](https://docs.docker.com/engine/installation) (optional)
  * [Docker compose](https://docs.docker.com/compose/install/) (optional)


### Run development in docker

There is docker-compose file allow to start docker container for development purpose.
This container shared source code dir with host.

To start container run this command

```
docker-compose up -d
```

Once the docker-compose environment is running, you can attach to the container with this command

```
docker exec -it github-authorized-keys sh
```

Source code is shared into ``/go/src/github.com/cloudposse/github-authorized-keys`` directory.

**Install dev tools inside container**

```
apk update
apk add git make curl
curl https://glide.sh/get | sh
```


### Install go libs dependencies

  Run ``make deps-dev`` to install additional go libs

### Testing

**Warning:**
Tests require perms to create users, so it is better to run them inside docker container.

Running tests required some configs.
There are 2 approaches to do this.

#### With config file

Copy .github-authorized-keys-tests.default.yml to .github-authorized-keys-tests.yml

```
cp .github-authorized-keys-tests.default.yml .github-authorized-keys-tests.yml
```

and set required values in that file.

Then you can simple run

```
make test
```

#### With environment variables

Run tests with command


```
TEST_GITHUB_API_TOKEN={api token} \
TEST_GITHUB_ORGANIZATION={organization name} \
TEST_GITHUB_TEAM={team name} \
TEST_GITHUB_TEAM_ID={team id} \
TEST_GITHUB_USER={user} \
TEST_ETCDCTL_ENDPOINT={etcd endpoints comma separeted - optional}
make test
```


### Run tests on docker build

To enable test run on docker build use ``--build-arg`` option
to set ``RUN_TESTS=1``

**Example**

```
docker build --build-arg RUN_TESTS=1 ./
```

Also you need to config tests before build
There are 2 approaches to do this.

#### With config file

The same way as described in configuration of tests with config file

#### With build args

Pass tests config environment variables as build-args

**Example**

```
docker build \
--build-arg RUN_TESTS=1 \
--build-arg  TEST_GITHUB_API_TOKEN={token} \
--build-arg  TEST_GITHUB_ORGANIZATION={org} \
--build-arg  TEST_GITHUB_TEAM={team} \
--build-arg  TEST_GITHUB_TEAM_ID={team_id} \
--build-arg  TEST_GITHUB_USER={user}
--build-arg  TEST_ETCDCTL_ENDPOINT={etcd endpoints comma separeted - optional}
```