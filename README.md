# Github Authorized Keys [![Build Status](https://travis-ci.org/cloudposse/github-authorized-keys.svg?branch=master)](https://travis-ci.org/cloudposse/github-authorized-keys)

Use GitHub teams to manage system user accounts and `authorized_keys`.

[![Docker Stars](https://img.shields.io/docker/stars/cloudposse/github-authorized-keys.svg)](https://hub.docker.com/r/cloudposse/github-authorized-keys)
[![Docker Pulls](https://img.shields.io/docker/pulls/cloudposse/github-authorized-keys.svg)](https://hub.docker.com/r/cloudposse/github-authorized-keys)
[![GitHub Stars](https://img.shields.io/github/stars/cloudposse/github-authorized-keys.svg)](https://github.com/cloudposse/github-authorized-keys/stargazers) 
[![GitHub Issues](https://img.shields.io/github/issues/cloudposse/github-authorized-keys.svg)](https://github.com/cloudposse/github-authorized-keys/issues)
[![Average time to resolve an issue](http://isitmaintained.com/badge/resolution/cloudposse/github-authorized-keys.svg)](http://isitmaintained.com/project/cloudposse/github-authorized-keys "Average time to resolve an issue")
[![Percentage of issues still open](http://isitmaintained.com/badge/open/cloudposse/github-authorized-keys.svg)](http://isitmaintained.com/project/cloudposse/github-authorized-keys "Percentage of issues still open")
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](https://github.com/cloudposse/github-authorized-keys/pulls)
[![License](https://img.shields.io/badge/license-APACHE%202.0%20-brightgreen.svg)](https://github.com/cloudposse/github-authorized-keys/blob/master/LICENSE)

----

## Architecture

This tool consists of three parts:

1. User Account / Authorized Keys provisioner which polls GitHub API for users that correspond to a given GitHub Organization & Team. It's responsible for adding or removing users from the system. All commands are templatized to allow it to run on multiple distributions. 
2. Simple read-only REST API that provides public keys for users, which is used by the `AuthorizedKeysCommand` in the `sshd_config`; this allows you to expose the service internally without compromising your Github Token. The public SSH access keys are *optionally* cached in Etcd for performance and reliability.
3. An `AuthorizedKeysCommand` [script](contrib/authorized-keys) that will `curl` the REST API for a user's public keys.

## Getting Started

By far, the easiest way to get up and running is by using the ready-made docker container. The only dependency is [Docker](https://docs.docker.com/engine/installation) itself.

We provide a stable public image [cloudposse/github-authorized-keys](https://hub.docker.com/r/cloudposse/github-authorized-keys/) or you can build your own from source.

```
docker build -t cloudposse/github-authorized-keys .
```

### Running GitHub Authorized Keys

All arguments can be passed both as environment variables or command-line arguments, or even mix-and-matched.

Available configuration options:

| **Environment Variable** | **Argument**             | **Description**                            | **Default** |
|--------------------------|--------------------------|--------------------------------------------|
| `GITHUB_API_TOKEN`       | `--github-api-token`     | GitHub API Token (read-only)               |
| `GITHUB_ORGANIZATION`    | `--github-organization`  | GitHub Organization Containing Team        |
| `GITHUB_TEAM`            | `--github-team`          | GitHub Team Membership to Grant SSH Access |
| `SYNC_USERS_GID`         | `--sync-users-gid`       | Default Group ID (aka `gid`) of users      |
| `SYNC_USERS_GROUPS`      | `--sync-users-groups`    | Default "Extra" Groups                     |
| `SYNC_USERS_SHELL`       | `--sync-users-shell`     | Default Login Shell                        |
| `SYNC_USERS_ROOT`        | `--sync-users-root`      | `chroot` path for user commands            |
| `SYNC_USERS_INTERVAL`    | `--sync-users-interval`  | Interval used to update user accounts      |
| `ETCD_ENDPOINT`          | `--etcd-endpoint`        | Etcd endpoint used for caching public keys |
| `ETCD_TTL`               | `--etcd-ttl`             | Duration (in seconds) to cache public keys |
| `ETCD_PREFIX`            | `--etcd-prefix`          | Prefix for public keys stored in etcd      |
| `LISTEN`                 | `--listen`               | Bind address used for REST API             |
| `INTEGRATE_SSH`          | `--integrate-ssh`        | Flag to automatically configure SSH        | `false`
| `LOG_LEVEL`              | `--log-level`            | Ccontrol the logging verbosity.            | `info`     |

## Quick Start 

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
  -e ETCD_ENDPOINT={etcd endpoints comma separeted - optional} \
  -e ETCD_TTL={etcd ttl - default 1 day} \
  -e ETCD_PREFIX={prefix or path to store data - default /github-authorized-keys} \
  -e LISTEN={Sets the address and port for IP, default :301} \
  -e INTEGRATE_SSH={integrate with ssh on startup, default false (should be true for production)} \
     cloudposse/github-authorized-keys:latest
```

## Usage Examples

### Automatically Configure SSH

To leverage the `github-authorized-keys` API, we need to make a small tweak to the `sshd_config`. 

This can be done automatically by passing the `--integrate-ssh` flag (or setting `INTEGRATE_SSH=true`)

After modifying the `sshd_config`, it's necessary to restart the SSH daemon. This happens automatically by calling the `SSH_RESTART_TPL` command. Since this differs depending on the OS distribution, you can change the default behavior by setting the `SSH_RESTART_TPL` environment variable (default: `/usr/sbin/service ssh force-reload`). Similarly, you might need to tweak the `AUTHORIZED_KEYS_COMMAND_TPL` environment variable to something compatible with your OS.


### Manually Configure SSH

If you wish to manually configure your `sshd_config`, here's all you need to do:

```
AuthorizedKeysCommand /usr/bin/authorized-keys
AuthorizedKeysCommandUser root
```

Then install a [wrapper script](contrib/authorized-keys) to `/usr/bin/authorized-keys`. 

**Note**: this command requires `curl` to access the REST API in order to fetch authorized keys

### Etcd Fallback Cache

Authorization REST API use ETCD to temporary cache user's public keys.
If github.com is not available command fallback to ETCD storage.

Etcd endpoints param is optional, if not specify caching and fallback disabled.

### Create users

Linux users will be synchronized according to team members every 1 second

In case of running in container you have to share host ``/`` into ``/{root directory}`` because  ``adduser`` command could differs on different Linux distribs and we need to use host one.
Also that means you need to specify  sync-users-root param to point to that directory.

### Templating Commands

Due to the differences between OS commands, the defaults might not work. 

Below are some of the settings which can be tweaked. 

| Environment Variable           | **Description **                                                                | **Default**                                   
|--------------------------------|---------------------------------------------------------------------------------|-------------------------------------------------------------------------
| `LINUX_USER_ADD_TPL`           | Command used to add a user to the system when no default group supplied.        | `adduser {username} --disabled-password --force-badname --shell {shell}`                 
| `LINUX_USER_ADD_WITH_GID_TPL`  | Command used to add a user to the system when a default primary group supplied. | `adduser {username} --disabled-password --force-badname --shell {shell} --group {group}`
| `LINUX_USER_ADD_TO_GROUP_TPL`  | Command used to add the user to secondary groups                                | `adduser {username} {group}` 
| `LINUX_USER_DEL_TPL`           | Command used to delete a user from the system when removed the the team         | `deluser {username}`
| `SSH_RESTART_TPL`              | Command used to restart SSH when `INTEGRATE_SSH=true`                           | `/usr/sbin/service ssh force-reload`

**Macros:**

1. `{username}` - User login name
2. `{shell}`    - User shell
3. `{group}`    - User primary group name
4. `{gid}`      - User primary group id

## Help

**Got a question?** 

File a GitHub [issue](https://github.com/cloudposse/github-authorized-keys/issues), send us an [email](mailto:hello@cloudposse.com) or reach out to us on [Gitter](https://gitter.im/cloudposse/).

## Contributing

### Bug Reports & Feature Requests

Please use the [issue tracker](https://github.com/cloudposse/github-authorized-keys/issues) to report any bugs or file feature requests.

### Developing

If you are interested in being a contributor and want to get involved in developing Geodesic, we would love to hear from you! Shoot us an [email](mailto:hello@cloudposse.com).

In general, PRs are welcome. We follow the typical "fork-and-pull" Git workflow.

 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull request** so that we can review your changes

**NOTE:** Be sure to merge the latest from "upstream" before making a pull request!

Here's how to get started...

1. `git clone https://github.com/cloudposse/github-authorized-keys.git` to pull down the repository 
2. `make init` to initialize the [`build-harness`](https://github.com/cloudposse/build-harness/)
3. Review the [documentation](docs/) on compiling

## License

[APACHE 2.0](LICENSE) © 2016-2017 [Cloud Posse, LLC](https://cloudposse.com)

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at
     
      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.

## About

GitHub Authorized Keys is maintained and funded by [Cloud Posse, LLC][website]. Like it? Please let us know at <hello@cloudposse.com>

We love [Open Source Software](https://github.com/cloudposse/)! 

See [our other projects][community] or [hire us][hire] to help build your next cloud-platform.

  [website]: http://cloudposse.com/
  [community]: https://github.com/cloudposse/
  [hire]: http://cloudposse.com/contact/
  
### Contributors


| [![Erik Osterman][erik_img]][erik_web]<br/>[Erik Osterman][erik_web] | [![Igor Rodionov][igor_img]][igor_web]<br/>[Igor Rodionov][igor_web] |
|-------------------------------------------------------|------------------------------------------------------------------|

  [erik_img]: http://s.gravatar.com/avatar/88c480d4f73b813904e00a5695a454cb?s=144
  [erik_web]: https://github.com/osterman/
  [igor_img]: http://s.gravatar.com/avatar/bc70834d32ed4517568a1feb0b9be7e2?s=144
  [igor_web]: https://github.com/goruha/


