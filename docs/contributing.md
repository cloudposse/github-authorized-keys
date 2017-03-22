# Compiling from Source / Development

## Requirements

  * [Go lang 1.7.x](https://golang.org/)
  * [glide](https://github.com/Masterminds/glide)
  * [Make](https://en.wikipedia.org/wiki/Make_(software))
  * [Docker](https://docs.docker.com/engine/installation) (optional)
  * [Docker compose](https://docs.docker.com/compose/install/) (optional)


## Quick Start

Run the following commands:
```
make go:deps
make go:build
make go:install
```

After installation, the binary will be installed here:
```
/usr/local/sbin/github-authorized-keys
```


## Use Docker Compose for Development

There is `docker-compose` file which starts a docker container for development purposes. This container shares source code directory from the host.

To start the container, run this command:

```
docker-compose up -d
```

Once the `docker-compose` environment is running, you can attach to the container with this command:

```
docker exec -it github-authorized-keys sh
```

Source code is bind-mounted into the `/go/src/github.com/cloudposse/github-authorized-keys` directory.

**Install developement tools inside of the container**

```
apk update
apk add git make curl
curl https://glide.sh/get | sh
```

## Install Golang Dependencies

Run `make go:deps-dev` to install additional go libs.

## Testing

**Warning:** Tests require sufficient permission to create users, therefore we recommend that you run them inside of a container

Running tests requires configuration. There are 2 approaches to do this:

### Using Config File

Copy `.github-authorized-keys-tests.default.yml` to `.github-authorized-keys-tests.yml`.

```
cp .github-authorized-keys-tests.default.yml .github-authorized-keys-tests.yml
```

Then update the settings inside of the `.github-authorized-keys-tests.yml` file.

After that, simply run:

```
make go:test
```

### Using Environment Variables

Run tests with command


```
TEST_GITHUB_API_TOKEN={api token} \
TEST_GITHUB_ORGANIZATION={organization name} \
TEST_GITHUB_TEAM={team name} \
TEST_GITHUB_TEAM_ID={team id} \
TEST_GITHUB_USER={user} \
TEST_ETCD_ENDPOINT={etcd endpoints comma separeted - optional} \
  make test
```


## Run Integration Tests on Docker Build

To enable test run on docker build use `--build-arg` option to set `RUN_TESTS=1`

**Example**

```
docker build --build-arg RUN_TESTS=1 ./
```

**Note:** You need to config tests before the `build` step. There are two ways to go about this:

1. The same way as described in configuration of tests with config file
2. Pass tests config environment variables as build-args

**Example**

```
docker build \
  --build-arg RUN_TESTS=1 \
  --build-arg  TEST_GITHUB_API_TOKEN={token} \
  --build-arg  TEST_GITHUB_ORGANIZATION={org} \
  --build-arg  TEST_GITHUB_TEAM={team} \
  --build-arg  TEST_GITHUB_TEAM_ID={team_id} \
  --build-arg  TEST_GITHUB_USER={user} \
  --build-arg  TEST_ETCD_ENDPOINT={etcd endpoints comma separeted - optional}
    ./
```
