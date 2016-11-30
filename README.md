# Github Authorized Keys
Allow to provide ssh access to servers based on github teams

----

## Demo

We use Vagrant to demonstrate how this tool works.

### Deps

Install

**Virtual box** (tested on version 4.3.26) https://www.virtualbox.org/wiki/Downloads

**Vagrant** (tested on version 1.8.4) https://www.vagrantup.com/downloads.html

**vagrant-docker-compose** plugin  with command

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

``ssh -o "UserKnownHostsFile /dev/null" {user}@192.168.33.10``

### Logs

You can check what is going with ssh inside vagrant box

``sudo tail -f /var/log/auth.log``

----------------

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

------------

## Development

### Requirements

  * Go lang 1.7.x
  * Make
  * Docker (optional)


### Run development in docker

There is docker-compose file allow to start docker container for development purpose.
This container shared source code dir with host.

To ssh inside container run

```
docker exec -it github-authorized-keys sh
cd /go/src/github.com/cloudposse/github-authorized-keys
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
docker build  \
--build-arg RUN_TESTS=1 \
--build-arg  TEST_GITHUB_API_TOKEN={token} \
--build-arg  TEST_GITHUB_ORGANIZATION={org} \
--build-arg  TEST_GITHUB_TEAM={team} \
--build-arg  TEST_GITHUB_TEAM_ID={team_id} \
--build-arg  TEST_GITHUB_USER={user}
./
```
